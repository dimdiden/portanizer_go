package app

import (
	"errors"
)

var (
	// ErrNotFound is an implementation agnostic error that should be returned
	// by any service implementation when a record was not located.
	ErrNotFound = errors.New("Record not found")
)

type Assigner interface {
	PutTags(pid string, tagids []string) (*Post, error)
}

type PostStore interface {
	Assigner
	GetByID(id string) (*Post, error)
	GetByName(name string) (*Post, error)
	GetList() ([]*Post, error)
	Create(post Post) (*Post, error)
	Update(id string, post Post) (*Post, error)
	Delete(id string) error
}

type TagStore interface {
	GetByID(id string) (*Tag, error)
	GetByName(name string) (*Tag, error)
	GetList() ([]*Tag, error)
	Create(tag Tag) (*Tag, error)
	Update(id string, tag Tag) (*Tag, error)
	Delete(id string) error
}
