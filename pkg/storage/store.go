package storage

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"slices"
)

type Store interface {
	Read(v any) error
	Write(v any) error
	Modified() bool
	Exist() bool
}

type FileStore struct {
	Path string
	stat os.FileInfo
}

func (s *FileStore) Read(v any) error {
	if _, err := os.Stat(s.Path); err != nil {
		return nil
	}

	bytes, err := os.ReadFile(s.Path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &v); err != nil {
		return err
	}
	return nil
}

func (s *FileStore) Write(v any) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.Path, bytes, os.ModePerm); err == nil {
		s.stat, _ = os.Stat(s.Path)
		return nil
	} else {
		return err
	}
}

func (s FileStore) Modified() bool {
	if s.stat == nil {
		return true
	}

	if stat, err := os.Stat(s.Path); err != nil {
		return false
	} else {
		return stat.ModTime().After(s.stat.ModTime())
	}
}

func (s FileStore) Exist() bool {
	_, err := os.Stat(s.Path)
	return err == nil
}

type HasId interface {
	GetId() string
	SetId(id string)
}

func GenerateId[T HasId](data []T) string {
	randomString := func(len int) (string, error) {
		b := make([]byte, len)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		return hex.EncodeToString(b), nil
	}

	for {
		if id, err := randomString(4); err != nil {
			continue
		} else if index, _ := IndexOf(id, data); index != -1 {
			continue
		} else {
			return id
		}
	}
}

func IndexOf[T HasId](id string, data []T) (int, error) {
	index := slices.IndexFunc(data, func(g T) bool {
		return g.GetId() == id
	})

	if index == -1 {
		return -1, errors.New("store: no item with the same ID was found")
	}
	return index, nil
}

func Create[T HasId](v T, data *[]T) (string, error) {
	v.SetId(GenerateId(*data))
	*data = append(*data, v)
	return v.GetId(), nil
}

func Update[T HasId](v T, data *[]T) error {
	if index, err := IndexOf(v.GetId(), *data); err != nil {
		return err
	} else {
		(*data)[index] = v
		return nil
	}
}

func Delete[T HasId](id string, data *[]T) error {
	if index, err := IndexOf(id, *data); err != nil {
		return err
	} else {
		*data = append((*data)[:index], (*data)[index+1:]...)
		return nil
	}
}

func Get[T HasId](id string, data []T) (T, error) {
	if index, err := IndexOf(id, data); err != nil {
		return *new(T), err
	} else {
		return data[index], nil
	}
}
