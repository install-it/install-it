package storage

import (
	"encoding/json"
	"errors"
	"os"
	"slices"
)

type DriverGroupManager struct {
	Path   string
	groups []DriverGroup
	fstat  os.FileInfo
}

func (m *DriverGroupManager) Read() ([]DriverGroup, error) {
	if m.fstat == nil || m.modified() {
		var groups []DriverGroup

		if _, err := os.Stat(m.Path); err != nil {
			os.WriteFile(m.Path, []byte("[]"), os.ModePerm)
		}

		bytes, err := os.ReadFile(m.Path)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(bytes, &groups); err != nil {
			return nil, err
		}

		m.groups = groups
	}
	return m.groups, nil
}

func (m *DriverGroupManager) write() error {
	bytes, err := json.Marshal(m.groups)
	if err != nil {
		return err
	}

	if err := os.WriteFile(m.Path, bytes, os.ModePerm); err == nil {
		m.fstat, _ = os.Stat(m.Path)
		return nil
	} else {
		return err
	}
}

func (m DriverGroupManager) modified() bool {
	if m.fstat == nil {
		return false
	}

	if stat, err := os.Stat(m.Path); err != nil {
		return false
	} else {
		return stat.ModTime().After(m.fstat.ModTime())
	}
}

func (m DriverGroupManager) IndexOf(groupId string) (int, error) {
	index := slices.IndexFunc(m.groups, func(g DriverGroup) bool {
		return g.Id == groupId
	})

	if index == -1 {
		return -1, errors.New("storage: no group with the same ID was found")
	}
	return index, nil
}

func (m DriverGroupManager) GroupOf(driverId string) (string, error) {
	for _, group := range m.groups {
		for _, driver := range group.Drivers {
			if driver.Id == driverId {
				return group.Id, nil
			}
		}
	}
	return "", errors.New("storage: no driver with the same ID was found in any group")
}

func (m *DriverGroupManager) Get(id string) (DriverGroup, error) {
	if index, err := m.IndexOf(id); err != nil {
		return DriverGroup{}, err
	} else {
		return m.groups[index], nil
	}
}

func (m DriverGroupManager) generateGid() string {
	for {
		if id, err := randomString(4); err != nil {
			continue
		} else if _, err := m.GroupOf(id); err == nil {
			continue
		} else {
			return id
		}
	}
}

func (m *DriverGroupManager) Add(group DriverGroup) (string, error) {
	for group.Id = ""; group.Id == ""; {
		if id, err := randomString(4); err != nil {
			continue
		} else if idx, _ := m.IndexOf(id); idx != -1 {
			continue
		} else {
			group.Id = id
		}
	}

	for idx := range group.Drivers {
		for group.Drivers[idx].Id = ""; group.Drivers[idx].Id == ""; {
			group.Drivers[idx].Id = m.generateGid()
		}
	}

	m.groups = append(m.groups, group)
	return group.Id, m.write()
}

func (m *DriverGroupManager) Update(group DriverGroup) error {
	if index, err := m.IndexOf(group.Id); err != nil {
		return err
	} else {
		for idx := range group.Drivers {
			if group.Drivers[idx].Id == "" {
				group.Drivers[idx].Id = m.generateGid()
			}
		}

		m.groups[index] = group
		return m.write()
	}
}

func (m *DriverGroupManager) Remove(id string) error {
	if index, err := m.IndexOf(id); err != nil {
		return err
	} else {
		// for _, group := range m.groups {
		// 	for _, driver := range group.Drivers {

		// 	}
		// }
		m.groups = append(m.groups[:index], m.groups[index+1:]...)
		return m.write()
	}
}

func (m *DriverGroupManager) MoveBehind(id string, index int) ([]DriverGroup, error) {
	if srcIndex, err := m.IndexOf(id); err != nil {
		return m.groups, err
	} else {
		if index < -1 || index >= len(m.groups)-1 {
			return m.groups, errors.New("storage: target index out of bound")
		}

		if len(m.groups) == 1 || srcIndex-index == 1 {
			return m.groups, nil
		}

		if srcIndex <= index {
			for i := srcIndex; i < index+1; i++ {
				m.groups[i], m.groups[i+1] = m.groups[i+1], m.groups[i]
			}
		} else {
			for i := srcIndex; i > index+1; i-- {
				m.groups[i-1], m.groups[i] = m.groups[i], m.groups[i-1]
			}
		}
		return m.groups, m.write()
	}
}

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
