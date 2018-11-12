package gorm

import (
	"github.com/dimdiden/portanizer_go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type userRepo struct {
	DB *gorm.DB
}

// NewUserRepo returns user repository implementation in gorm
func NewUserRepo(db *gorm.DB) portanizer.UserRepo {
	return &userRepo{DB: db}
}

func (r *userRepo) GetByCreds(email, pwd string) (*portanizer.User, error) {
	var user portanizer.User
	if r.DB.Where("email = ?", email).First(&user).RecordNotFound() {
		return nil, portanizer.ErrNotFound
	}
	if err := comparePasswords(user.Password, pwd); err != nil {
		return nil, portanizer.ErrAuth
	}
	return &user, nil
}

func (r *userRepo) Create(user portanizer.User) (*portanizer.User, error) {
	if !r.DB.First(&user, "email = ?", user.Email).RecordNotFound() {
		return nil, portanizer.ErrExists
	}
	hash, err := hashAndSalt(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = string(hash)

	if err := r.DB.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Refresh(user *portanizer.User) error {
	if err := r.DB.Model(&user).Update(portanizer.User{RToken: user.RToken}).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepo) EmptyRToken(id string) error {
	var user portanizer.User
	if r.DB.First(&user, "id = ?", id).RecordNotFound() {
		return portanizer.ErrNotFound
	}
	if err := r.DB.Model(&user).Update("r_token", "").Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepo) GetByID(id string) (*portanizer.User, error) {
	var user portanizer.User
	if r.DB.First(&user, "id = ?", id).RecordNotFound() {
		return nil, portanizer.ErrNotFound
	}
	return &user, nil
}

func hashAndSalt(pwd string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func comparePasswords(hashPwd string, plainPwd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(plainPwd))
	if err != nil {
		return err
	}
	return nil
}
