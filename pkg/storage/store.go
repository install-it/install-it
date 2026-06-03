package storage

import (
	"encoding/json"
	"os"
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

	if err := os.WriteFile(s.Path, bytes, 0777); err != nil {
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
