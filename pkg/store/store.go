package store

import (
	"encoding/json"
	"os"
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
