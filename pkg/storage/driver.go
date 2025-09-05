package storage

import (
	"driver-box/pkg/utils"
	"errors"
	"slices"
)

type DriverGroup struct {
	Id      string     `json:"id"`
	Name    string     `json:"name"`
	Type    DriverType `json:"type"`
	Drivers []*Driver  `json:"drivers"`
}

func (g DriverGroup) GetId() string { return g.Id }

func (g *DriverGroup) SetId(id string) { g.Id = id }

type DriverType string

const (
	Network       DriverType = "network"
	Display       DriverType = "display"
	Miscellaneous DriverType = "miscellaneous"
)

type Driver struct {
	Id            string     `json:"id"`
	Name          string     `json:"name"`
	Type          DriverType `json:"type"`
	Path          string     `json:"path"`
	Flags         []string   `json:"flags"`
	MinExeTime    float32    `json:"minExeTime"`
	AllowRtCodes  []int32    `json:"allowRtCodes"`
	Incompatibles []string   `json:"incompatibles"`
}

func (d Driver) GetId() string { return d.Id }

func (d *Driver) SetId(id string) { d.Id = id }

type DriverGroupStorage struct {
	Store Store
	data  []*DriverGroup
}

func (s *DriverGroupStorage) All() ([]*DriverGroup, error) {
	if s.Store.Modified() {
		s.data = []*DriverGroup{}
		s.Store.Read(&s.data)
	} else if !s.Store.Exist() {
		s.data = []*DriverGroup{}
		s.Store.Write(s.data)
	}
	return s.data, nil
}

func (s *DriverGroupStorage) Get(id string) (DriverGroup, error) {
	if group, err := Get(id, s.data); err != nil {
		return DriverGroup{}, err
	} else {
		return *group, nil
	}
}

func (s *DriverGroupStorage) Add(group DriverGroup) (string, error) {
	var drivers []*Driver
	for _, g := range s.data {
		drivers = append(drivers, g.Drivers...)
	}

	for i := range group.Drivers {
		group.Drivers[i].Id = GenerateId(drivers)
		drivers = append(drivers, group.Drivers[i])
	}

	if id, err := Create(&group, &s.data); err != nil {
		return "", err
	} else {
		return id, s.Store.Write(s.data)
	}
}

// func (s *DriverGroupStorage) Update(group DriverGroup) error {
// 	var drivers []*Driver
// 	var driverIds []string
// 	for _, g := range s.data {
// 		drivers = append(drivers, g.Drivers...)
// 		for _, d := range g.Drivers {
// 			driverIds = append(driverIds, d.Id)
// 		}
// 	}

// 	for i := range group.Drivers {
// 		if !slices.Contains(driverIds, group.Drivers[i].Id) {
// 			group.Drivers[i].Id = GenerateId(drivers)
// 			drivers = append(drivers, group.Drivers[i])
// 		}
// 	}

//		if err := Update(&group, &s.data); err != nil {
//			return err
//		}
//		return s.Store.Write(s.data)
//	}

func (s *DriverGroupStorage) Update(group DriverGroup) error {
	existingGroup, err := Get(group.Id, s.data)
	if err != nil {
		return err
	}

	// Find removed driver IDs
	driverMap := make(map[string]bool)
	for _, driver := range group.Drivers {
		driverMap[driver.Id] = true
	}

	removedDriverIds := make([]string, 0)
	for _, driver := range existingGroup.Drivers {
		if !driverMap[driver.Id] {
			removedDriverIds = append(removedDriverIds, driver.Id)
		}
	}

	// Generate IDs for new drivers
	existingIds := make(map[string]bool)
	for _, g := range s.data {
		for _, d := range g.Drivers {
			existingIds[d.Id] = true
		}
	}

	for i, driver := range group.Drivers {
		if !existingIds[driver.Id] {
			group.Drivers[i].Id = GenerateId([]*Driver{})
		}
	}

	// Update group and clean up incompatibles
	if err := Update(&group, &s.data); err != nil {
		return err
	}

	if len(removedDriverIds) > 0 {
		removedSet := make(map[string]bool)
		for _, id := range removedDriverIds {
			removedSet[id] = true
		}

		for _, g := range s.data {
			for _, d := range g.Drivers {
				d.Incompatibles = slices.DeleteFunc(d.Incompatibles, func(id string) bool {
					return removedSet[id]
				})
			}
		}
	}

	return s.Store.Write(s.data)
}

func (s *DriverGroupStorage) Update2(group DriverGroup) error {
	group, err := s.Get(group.Id)
	if err != nil {
		return err
	}

	var drivers []*Driver
	for _, g := range s.data {
		drivers = append(drivers, g.Drivers...)
	}

	driverIds := utils.Map(drivers, func(d *Driver) string { return d.Id })

	// generate ID for new drivers
	for i := range group.Drivers {
		if !slices.Contains(driverIds, group.Drivers[i].Id) {
			group.Drivers[i].Id = GenerateId(drivers)
			drivers = append(drivers, group.Drivers[i])
		}
	}

	// update group
	if err := Update(&group, &s.data); err != nil {
		return err
	}

	// cacased deletion on Driver.Incompatibles
	deletedIds := slices.DeleteFunc(driverIds, func(id string) bool {
		for _, d := range group.Drivers {
			if d.Id == id {
				return true
			}
		}
		return false
	})

	for _, g := range s.data {
		for _, d := range g.Drivers {
			d.Incompatibles = slices.DeleteFunc(d.Incompatibles, func(id string) bool {
				return slices.Contains(deletedIds, id)
			})
		}
	}

	return s.Store.Write(s.data)
}

func (s *DriverGroupStorage) Remove(id string) error {
	group, err := Get(id, s.data)
	if err != nil {
		return err
	}

	driverIds := make([]string, 0, len(group.Drivers))
	for _, driver := range group.Drivers {
		driverIds = append(driverIds, driver.Id)
	}

	if err := Delete(id, &s.data); err != nil {
		return err
	}

	for _, group := range s.data {
		for _, driver := range group.Drivers {
			driver.Incompatibles = slices.DeleteFunc(driver.Incompatibles, func(id string) bool {
				return slices.Contains(driverIds, id)
			})
		}
	}

	return s.Store.Write(s.data)
}

func (s DriverGroupStorage) IndexOf(id string) (int, error) {
	return IndexOf(id, s.data)
}

func (s *DriverGroupStorage) MoveBehind(id string, index int) ([]*DriverGroup, error) {
	if srcIndex, err := IndexOf(id, s.data); err != nil {
		return s.data, err
	} else {
		if index < -1 || index >= len(s.data)-1 {
			return s.data, errors.New("store: target index out of bound")
		}

		if len(s.data) == 1 || srcIndex-index == 1 {
			return s.data, nil
		}

		if srcIndex <= index {
			for i := srcIndex; i < index+1; i++ {
				s.data[i], s.data[i+1] = s.data[i+1], s.data[i]
			}
		} else {
			for i := srcIndex; i > index+1; i-- {
				s.data[i-1], s.data[i] = s.data[i], s.data[i-1]
			}
		}
		return s.data, s.Store.Write(s.data)
	}
}
