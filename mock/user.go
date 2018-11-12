package mock

import (
	"github.com/dimdiden/portanizer_go"
)

type UserRepo struct {
	GetByIDFn      func(id string) (*portanizer.User, error)
	GetByIDInvoked bool

	GetByCredsFn      func(email, pwd string) (*portanizer.User, error)
	GetByCredsInvoked bool

	CreateFn      func(user portanizer.User) (*portanizer.User, error)
	CreateInvoked bool

	EmptyRTokenFn      func(id string) error
	EmptyRTokenInvoked bool

	RefreshFn      func(user *portanizer.User) error
	RefreshInvoked bool
}

func (r *UserRepo) GetByID(id string) (*portanizer.User, error) {
	r.GetByIDInvoked = true
	return r.GetByIDFn(id)
}

func (r *UserRepo) GetByCreds(email, pwd string) (*portanizer.User, error) {
	r.GetByCredsInvoked = true
	return r.GetByCredsFn(email, pwd)
}

func (r *UserRepo) Create(user portanizer.User) (*portanizer.User, error) {
	r.CreateInvoked = true
	return r.CreateFn(user)
}

func (r *UserRepo) EmptyRToken(id string) error {
	r.EmptyRTokenInvoked = true
	return r.EmptyRTokenFn(id)
}

func (r *UserRepo) Refresh(user *portanizer.User) error {
	r.RefreshInvoked = true
	return r.RefreshFn(user)
}
