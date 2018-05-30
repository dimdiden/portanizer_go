package app

import "errors"

var (
	// ErrNotFound is an implementation agnostic error that should be returned
	// by any service implementation when a record was not located.
	ErrNotFound = errors.New("Record not found")
)

type TagStore interface {
	GetByID(id string) (*Tag, error)
	GetByName(name string) (*Tag, error)
	GetList() ([]*Tag, error)
	Create(tag Tag) (*Tag, error)
	Update(id string, tag Tag) (*Tag, error)
	Delete(id string) error
}
