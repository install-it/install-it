package store

import (
	"errors"
	"slices"
)

type DriverGroup struct {
	Id      string     `json:"id"`
	Name    string     `json:"name"`
	Type    DriverType `json:"type"`
	Drivers []Driver   `json:"drivers"`
}

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

type DriverGroupStore struct {
	Store Store
	data  []DriverGroup
}

func (s *DriverGroupStore) All() ([]DriverGroup, error) {
	if s.Store.Modified() {
		s.data = []DriverGroup{}
		s.Store.Read(&s.data)
	} else if !s.Store.Exist() {
		s.data = []DriverGroup{}
		s.Store.Write(s.data)
	}
	return s.data, nil
}

func (s *DriverGroupStore) Get(id string) (DriverGroup, error) {
	if index, err := s.IndexOf(id); err != nil {
		return DriverGroup{}, err
	} else {
		return s.data[index], nil
	}
}

func (s *DriverGroupStore) Add(group DriverGroup) (string, error) {
	for group.Id = ""; group.Id == ""; {
		if id, err := randomString(4); err != nil {
			continue
		} else if idx, _ := s.IndexOf(id); idx != -1 {
			continue
		} else {
			group.Id = id
		}
	}

	for idx := range group.Drivers {
		for group.Drivers[idx].Id = ""; group.Drivers[idx].Id == ""; {
			group.Drivers[idx].Id = s.generateGid()
		}
	}

	s.data = append(s.data, group)
	return group.Id, s.Store.Write(s.data)
}

func (s *DriverGroupStore) Update(group DriverGroup) error {
	if index, err := s.IndexOf(group.Id); err != nil {
		return err
	} else {
		for idx := range group.Drivers {
			if group.Drivers[idx].Id == "" {
				group.Drivers[idx].Id = s.generateGid()
			}
		}

		s.data[index] = group
		return s.Store.Write(s.data)
	}
}

func (s *DriverGroupStore) Remove(id string) error {
	if index, err := s.IndexOf(id); err != nil {
		return err
	} else {
		s.data = append(s.data[:index], s.data[index+1:]...)
		return s.Store.Write(s.data)
	}
}

func (s DriverGroupStore) IndexOf(id string) (int, error) {
	index := slices.IndexFunc(s.data, func(g DriverGroup) bool {
		return g.Id == id
	})

	if index == -1 {
		return -1, errors.New("store: no group with the same ID was found")
	}
	return index, nil
}

func (s DriverGroupStore) GroupOf(id string) (string, error) {
	for _, group := range s.data {
		for _, driver := range group.Drivers {
			if driver.Id == id {
				return group.Id, nil
			}
		}
	}
	return "", errors.New("store: no driver with the same ID was found in any group")
}

func (s *DriverGroupStore) MoveBehind(id string, index int) ([]DriverGroup, error) {
	if srcIndex, err := s.IndexOf(id); err != nil {
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

func (s DriverGroupStore) generateGid() string {
	for {
		if id, err := randomString(4); err != nil {
			continue
		} else if _, err := s.GroupOf(id); err == nil {
			continue
		} else {
			return id
		}
	}
}
