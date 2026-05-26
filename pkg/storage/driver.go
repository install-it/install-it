package storage

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"slices"

	"gorm.io/gorm"
)

type DriverType string

const (
	Network       DriverType = "network"
	Display       DriverType = "display"
	Miscellaneous DriverType = "miscellaneous"
)

type DriverGroup struct {
	Id                string    `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name"`
	Type              DriverType `json:"type"`
	MutuallyExclusive bool      `json:"mutuallyExclusive"`
	Position          int       `json:"-" gorm:"index"`
	Drivers           []*Driver `json:"drivers" gorm:"foreignKey:GroupId;constraint:OnDelete:CASCADE"`
}

func (g *DriverGroup) BeforeCreate(tx *gorm.DB) error {
	if g.Id == "" {
		g.Id = generateHexId()
	}
	return nil
}

type Driver struct {
	Id            string     `json:"id" gorm:"primaryKey"`
	GroupId       string     `json:"-" gorm:"index"`
	Name          string     `json:"name"`
	Type          DriverType `json:"type"`
	Path          string     `json:"path"`
	Flags         []string   `json:"flags" gorm:"serializer:json"`
	MinExeTime    float32    `json:"minExeTime"`
	AllowRtCodes  []int32    `json:"allowRtCodes" gorm:"serializer:json"`
	Incompatibles []string   `json:"incompatibles" gorm:"serializer:json"`
}

func (d *Driver) BeforeCreate(tx *gorm.DB) error {
	if d.Id == "" {
		d.Id = generateHexId()
	}
	return nil
}

func generateHexId() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

type DriverGroupStorage struct {
	DB *gorm.DB
}

func NewDriverGroupStorage(db *gorm.DB) *DriverGroupStorage {
	return &DriverGroupStorage{DB: db}
}

func (s *DriverGroupStorage) All() ([]DriverGroup, error) {
	var groups []*DriverGroup
	if err := s.DB.Preload("Drivers").Order("position").Find(&groups).Error; err != nil {
		return nil, err
	}
	result := make([]DriverGroup, len(groups))
	for i, g := range groups {
		result[i] = *g
	}
	return result, nil
}

func (s *DriverGroupStorage) Get(id string) (DriverGroup, error) {
	var group DriverGroup
	result := s.DB.Preload("Drivers").First(&group, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return DriverGroup{}, errors.New("store: no item with the same ID was found")
		}
		return DriverGroup{}, result.Error
	}
	return group, nil
}

func (s *DriverGroupStorage) Add(group DriverGroup) (string, error) {
	var maxPos int
	s.DB.Model(&DriverGroup{}).Select("COALESCE(MAX(position), -1)").Scan(&maxPos)
	group.Position = maxPos + 1

	if err := s.DB.Create(&group).Error; err != nil {
		return "", err
	}
	return group.Id, nil
}

func (s *DriverGroupStorage) Update(group DriverGroup) (DriverGroup, error) {
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Determine which driver IDs have been removed
		var existing []*Driver
		if err := tx.Where("group_id = ?", group.Id).Find(&existing).Error; err != nil {
			return err
		}

		newDriverIds := make(map[string]bool)
		for _, d := range group.Drivers {
			if d.Id != "" {
				newDriverIds[d.Id] = true
			}
		}

		var deletedIds []string
		for _, d := range existing {
			if !newDriverIds[d.Id] {
				deletedIds = append(deletedIds, d.Id)
			}
		}

		// Delete removed drivers
		if len(deletedIds) > 0 {
			if err := tx.Delete(&Driver{}, "id IN ?", deletedIds).Error; err != nil {
				return err
			}
		}

		// Update group scalar fields
		if err := tx.Model(&DriverGroup{}).Where("id = ?", group.Id).Updates(map[string]any{
			"name":               group.Name,
			"type":               group.Type,
			"mutually_exclusive": group.MutuallyExclusive,
		}).Error; err != nil {
			return err
		}

		// Upsert drivers
		for _, d := range group.Drivers {
			d.GroupId = group.Id
			if d.Id == "" {
				if err := tx.Create(d).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Save(d).Error; err != nil {
					return err
				}
			}
		}

		// Clean up Incompatibles in all remaining drivers
		if len(deletedIds) > 0 {
			var allDrivers []*Driver
			if err := tx.Find(&allDrivers).Error; err != nil {
				return err
			}
			for _, d := range allDrivers {
				cleaned := slices.DeleteFunc(slices.Clone(d.Incompatibles), func(id string) bool {
					return slices.Contains(deletedIds, id)
				})
				if len(cleaned) != len(d.Incompatibles) {
					if err := tx.Model(d).Update("incompatibles", cleaned).Error; err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return DriverGroup{}, err
	}
	return s.Get(group.Id)
}

func (s *DriverGroupStorage) Remove(id string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Collect driver IDs before CASCADE removes them
		var drivers []*Driver
		if err := tx.Where("group_id = ?", id).Find(&drivers).Error; err != nil {
			return err
		}
		driverIds := make([]string, len(drivers))
		for i, d := range drivers {
			driverIds[i] = d.Id
		}

		// Delete group (CASCADE removes its drivers)
		result := tx.Delete(&DriverGroup{}, "id = ?", id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("store: no item with the same ID was found")
		}

		if len(driverIds) == 0 {
			return nil
		}

		// Clean up Incompatibles in all remaining drivers
		var allDrivers []*Driver
		if err := tx.Find(&allDrivers).Error; err != nil {
			return err
		}
		for _, d := range allDrivers {
			cleaned := slices.DeleteFunc(slices.Clone(d.Incompatibles), func(incompId string) bool {
				return slices.Contains(driverIds, incompId)
			})
			if len(cleaned) != len(d.Incompatibles) {
				d.Incompatibles = cleaned
				if err := tx.Save(d).Error; err != nil {
					return err
				}
			}
		}

		// Clean up DriverGroupIds in all RuleSets (remove the group being deleted)
		var ruleSets []*RuleSet
		if err := tx.Find(&ruleSets).Error; err != nil {
			return err
		}
		for _, rs := range ruleSets {
			cleaned := slices.DeleteFunc(slices.Clone(rs.DriverGroupIds), func(gid string) bool {
				return gid == id
			})
			if len(cleaned) != len(rs.DriverGroupIds) {
				rs.DriverGroupIds = cleaned
				if err := tx.Save(rs).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (s *DriverGroupStorage) IndexOf(id string) (int, error) {
	var groups []*DriverGroup
	if err := s.DB.Select("id").Order("position").Find(&groups).Error; err != nil {
		return -1, err
	}
	for i, g := range groups {
		if g.Id == id {
			return i, nil
		}
	}
	return -1, errors.New("store: no item with the same ID was found")
}

func (s *DriverGroupStorage) MoveBehind(id string, index int) ([]DriverGroup, error) {
	var groups []*DriverGroup
	if err := s.DB.Order("position").Find(&groups).Error; err != nil {
		return nil, err
	}

	srcIndex := -1
	for i, g := range groups {
		if g.Id == id {
			srcIndex = i
			break
		}
	}
	if srcIndex == -1 {
		all, _ := s.All()
		return all, errors.New("store: no item with the same ID was found")
	}

	if index < -1 || index >= len(groups)-1 {
		all, _ := s.All()
		return all, errors.New("store: target index out of bound")
	}

	if len(groups) == 1 || srcIndex-index == 1 {
		return s.All()
	}

	// Reorder in memory (same algorithm as original)
	if srcIndex <= index {
		for i := srcIndex; i < index+1; i++ {
			groups[i], groups[i+1] = groups[i+1], groups[i]
		}
	} else {
		for i := srcIndex; i > index+1; i-- {
			groups[i-1], groups[i] = groups[i], groups[i-1]
		}
	}

	// Persist new positions
	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		for i, g := range groups {
			if err := tx.Model(&DriverGroup{}).Where("id = ?", g.Id).Update("position", i).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return s.All()
}
