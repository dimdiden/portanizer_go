package mock

import (
	"github.com/dimdiden/portanizer_go"
)

type TagStore struct {
	GetIdFn      func(id string) (*portanizer.Tag, error)
	GetIdInvoked bool

	GetNameFn      func(name string) (*portanizer.Tag, error)
	GetNameInvoked bool

	GetListFn      func() ([]*portanizer.Tag, error)
	GetListInvoked bool

	CreateFn      func(tag portanizer.Tag) (*portanizer.Tag, error)
	CreateInvoked bool

	UpdateFn      func(id string, tag portanizer.Tag) (*portanizer.Tag, error)
	UpdateInvoked bool

	DeleteFn      func(id string) error
	DeleteInvoked bool
}

func (s *TagStore) GetByID(id string) (*portanizer.Tag, error) {
	s.GetIdInvoked = true
	return s.GetIdFn(id)
}

func (s *TagStore) GetByName(name string) (*portanizer.Tag, error) {
	s.GetNameInvoked = true
	return s.GetNameFn(name)
}

func (s *TagStore) GetList() ([]*portanizer.Tag, error) {
	s.GetListInvoked = true
	return s.GetListFn()
}

func (s *TagStore) Create(tag portanizer.Tag) (*portanizer.Tag, error) {
	s.CreateInvoked = true
	return s.CreateFn(tag)
}

func (s *TagStore) Update(id string, tag portanizer.Tag) (*portanizer.Tag, error) {
	s.UpdateInvoked = true
	return s.UpdateFn(id, tag)
}

func (s *TagStore) Delete(id string) error {
	s.DeleteInvoked = true
	return s.DeleteFn(id)
}
