package storage

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"slices"
	"sync"
)

type Store interface {
	Read(v any) error
	Write(v any) error
	Exist() bool
}

type FileStore struct {
	Path string
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

	if err := os.WriteFile(s.Path, bytes, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (s FileStore) Exist() bool {
	_, err := os.Stat(s.Path)
	return err == nil
}

// MemoryStore implements Store interface for in-memory storage with no persistent mechanism
type MemoryStore struct {
	data  any
	mutex sync.RWMutex
}

func (m *MemoryStore) Read(v any) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.data == nil {
		return nil
	}

	// Marshal and unmarshal to simulate file storage behavior
	bytes, err := json.Marshal(m.data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, v); err != nil {
		return err
	}
	return nil
}

func (m *MemoryStore) Write(v any) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Validate by marshaling
	if _, err := json.Marshal(v); err != nil {
		return err
	}

	m.data = v
	return nil
}

func (m *MemoryStore) Exist() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.data != nil
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

// DeleteEventBus manages event subscribers for delete operations.
type DeleteEventBus struct {
	subscribers map[string][]func(ids []string) error
	mutex       sync.RWMutex
}

func NewEventBus() *DeleteEventBus {
	return &DeleteEventBus{
		subscribers: make(map[string][]func(ids []string) error),
	}
}

func (d *DeleteEventBus) Subscribe(storage string, handler func(ids []string) error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.subscribers[storage] = append(d.subscribers[storage], handler)
}

func (d *DeleteEventBus) Publish(storage string, ids []string) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for _, handler := range d.subscribers[storage] {
		if err := handler(ids); err != nil {
			return err
		}
	}

	return nil
}
