package store

import (
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

type DriverGroupStore struct {
	Store Store
	data  []*DriverGroup
}

func (s *DriverGroupStore) All() ([]*DriverGroup, error) {
	if s.Store.Modified() {
		s.data = []*DriverGroup{}
		s.Store.Read(&s.data)
	} else if !s.Store.Exist() {
		s.data = []*DriverGroup{}
		s.Store.Write(s.data)
	}
	return s.data, nil
}

func (s *DriverGroupStore) Get(id string) (DriverGroup, error) {
	if group, err := Get(id, s.data); err != nil {
		return DriverGroup{}, err
	} else {
		return *group, nil
	}
}

func (s *DriverGroupStore) Add(group DriverGroup) (string, error) {
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

func (s *DriverGroupStore) Update(group DriverGroup) error {
	if err := Update(&group, &s.data); err != nil {
		return err
	}
	return s.Store.Write(s.data)
}

func (s *DriverGroupStore) Remove(id string) error {
	if err := Delete(id, &s.data); err != nil {
		return err
	}

	for i, group := range s.data {
		if index := slices.IndexFunc(group.Drivers, func(d *Driver) bool { return d.Id == id }); index != -1 {
			s.data[i].Drivers = append(group.Drivers[:index], group.Drivers[index+1:]...)
		}
	}

	return s.Store.Write(s.data)
}

func (s DriverGroupStore) IndexOf(id string) (int, error) {
	return IndexOf(id, s.data)
}

func (s *DriverGroupStore) MoveBehind(id string, index int) ([]*DriverGroup, error) {
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
