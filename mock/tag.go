package mock

import (
	"github.com/dimdiden/portanizer_go"
)

type TagRepo struct {
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

func (r *TagRepo) GetByID(id string) (*portanizer.Tag, error) {
	r.GetIdInvoked = true
	return r.GetIdFn(id)
}

func (r *TagRepo) GetByName(name string) (*portanizer.Tag, error) {
	r.GetNameInvoked = true
	return r.GetNameFn(name)
}

func (r *TagRepo) GetList() ([]*portanizer.Tag, error) {
	r.GetListInvoked = true
	return r.GetListFn()
}

func (r *TagRepo) Create(tag portanizer.Tag) (*portanizer.Tag, error) {
	r.CreateInvoked = true
	return r.CreateFn(tag)
}

func (r *TagRepo) Update(id string, tag portanizer.Tag) (*portanizer.Tag, error) {
	r.UpdateInvoked = true
	return r.UpdateFn(id, tag)
}

func (r *TagRepo) Delete(id string) error {
	r.DeleteInvoked = true
	return r.DeleteFn(id)
}
