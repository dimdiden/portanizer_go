package mock

import (
	app "github.com/dimdiden/portanizer_sop"
)

type TagStore struct {
	GetIdFn      func(id string) (*app.Tag, error)
	GetIdInvoked bool

	GetNameFn      func(name string) (*app.Tag, error)
	GetNameInvoked bool

	GetListFn      func() ([]*app.Tag, error)
	GetListInvoked bool

	CreateFn      func(tag app.Tag) (*app.Tag, error)
	CreateInvoked bool

	UpdateFn      func(id string, tag app.Tag) (*app.Tag, error)
	UpdateInvoked bool

	DeleteFn      func(id string) error
	DeleteInvoked bool
}

func (s *TagStore) GetByID(id string) (*app.Tag, error) {
	s.GetIdInvoked = true
	return s.GetIdFn(id)
}

func (s *TagStore) GetByName(name string) (*app.Tag, error) {
	s.GetNameInvoked = true
	return s.GetNameFn(name)
}

func (s *TagStore) GetList() ([]*app.Tag, error) {
	s.GetListInvoked = true
	return s.GetListFn()
}

func (s *TagStore) Create(tag app.Tag) (*app.Tag, error) {
	s.CreateInvoked = true
	return s.CreateFn(tag)
}

func (s *TagStore) Update(id string, tag app.Tag) (*app.Tag, error) {
	s.UpdateInvoked = true
	return s.UpdateFn(id, tag)
}

func (s *TagStore) Delete(id string) error {
	s.DeleteInvoked = true
	return s.DeleteFn(id)
}
