package storage

import (
	"errors"
	"install-it/pkg/utils"
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

func (s *DriverGroupStorage) All() ([]DriverGroup, error) {
	if !s.Store.Exist() {
		s.data = []*DriverGroup{}
		s.Store.Write(s.data)
	} else {
		s.Store.Read(&s.data)
	}
	return s.copyOfAll(), nil
}

func (s DriverGroupStorage) Get(id string) (DriverGroup, error) {
	if group, err := Get(id, s.data); err != nil {
		return DriverGroup{}, err
	} else {
		return *group, nil
	}
}

func (s *DriverGroupStorage) Add(group DriverGroup) (string, error) {
	drivers := utils.FlatMap(s.data, func(g *DriverGroup) []*Driver { return g.Drivers })

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

func (s *DriverGroupStorage) Update(group DriverGroup) (DriverGroup, error) {
	// slice of all the existing drivers
	drivers := utils.FlatMap(s.data, func(g *DriverGroup) []*Driver { return g.Drivers })
	// slice of all the existing drivers' ID
	driverIds := utils.Map(drivers, func(d *Driver) string { return d.Id })
	// slice of all the drivers' ID that will be delete after the update
	deletedIds := slices.DeleteFunc(driverIds, func(id string) bool {
		for _, d := range group.Drivers {
			if d.Id == id {
				return true
			}
		}
		return false
	})

	// generate ID for new drivers
	for i := range group.Drivers {
		if !slices.Contains(driverIds, group.Drivers[i].Id) {
			group.Drivers[i].Id = GenerateId(drivers)
			drivers = append(drivers, group.Drivers[i])
		}
	}

	// update group
	if err := Update(&group, &s.data); err != nil {
		return DriverGroup{}, err
	}

	// cacased deletion on Driver.Incompatibles
	for _, g := range s.data {
		for _, d := range g.Drivers {
			d.Incompatibles = slices.DeleteFunc(d.Incompatibles, func(id string) bool {
				return slices.Contains(deletedIds, id)
			})
		}
	}

	return group, s.Store.Write(s.data)
}

func (s *DriverGroupStorage) Remove(id string) error {
	group, err := s.Get(id)
	if err != nil {
		return err
	}
	driverIds := utils.Map(group.Drivers, func(d *Driver) string { return d.Id })

	if err := Delete(id, &s.data); err != nil {
		return err
	}

	for _, g := range s.data {
		for _, d := range g.Drivers {
			d.Incompatibles = slices.DeleteFunc(d.Incompatibles, func(id string) bool {
				return slices.Contains(driverIds, id)
			})
		}
	}

	return s.Store.Write(s.data)
}

func (s DriverGroupStorage) IndexOf(id string) (int, error) {
	return IndexOf(id, s.data)
}

func (s *DriverGroupStorage) MoveBehind(id string, index int) ([]DriverGroup, error) {
	if srcIndex, err := s.IndexOf(id); err != nil {
		return s.copyOfAll(), err
	} else {
		if index < -1 || index >= len(s.data)-1 {
			return s.copyOfAll(), errors.New("store: target index out of bound")
		}

		if len(s.data) == 1 || srcIndex-index == 1 {
			return s.copyOfAll(), nil
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

		return s.copyOfAll(), s.Store.Write(s.data)
	}
}

func (s DriverGroupStorage) copyOfAll() []DriverGroup {
	return utils.Map(s.data, func(g *DriverGroup) DriverGroup { return *g })
}
