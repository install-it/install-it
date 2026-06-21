package storage

import "errors"

// ErrNotFound is returned when a storage operation finds no matching record.
var ErrNotFound = errors.New("storage: not found")