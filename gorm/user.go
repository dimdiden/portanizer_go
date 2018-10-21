package gorm

import (
	"github.com/dimdiden/portanizer_go"
	"github.com/jinzhu/gorm"
)

type userRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) portanizer.UserRepo {
	return &userRepo{DB: db}
}

func (r *userRepo) Exists(user portanizer.User) bool {
	if r.DB.Where("email = ?", user.Email).First(&user).RecordNotFound() {
		return false
	}
	return true
}

func (r *userRepo) Create(user portanizer.User) (*portanizer.User, error) {
	if !r.DB.First(&user, "email = ?", user.Email).RecordNotFound() {
		return nil, portanizer.ErrExists
	}
	if err := r.DB.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
