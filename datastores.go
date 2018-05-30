package app

type TagStore interface {
	GetByID(id string) (*Tag, error)
	GetByName(name string) (*Tag, error)
}
