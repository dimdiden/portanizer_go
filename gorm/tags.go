package gorm

import (
	"fmt"

	"github.com/dimdiden/portanizer_go"
	"github.com/jinzhu/gorm"
)

type tagRepo struct {
	DB *gorm.DB
}

func NewTagRepo(db *gorm.DB) portanizer.TagRepo {
	return &tagRepo{DB: db}
}

func (r *tagRepo) GetByID(id string) (*portanizer.Tag, error) {
	var tag portanizer.Tag
	if r.DB.First(&tag, "id = ?", id).RecordNotFound() {
		return nil, portanizer.ErrNotFound
	}
	return &tag, nil
}

func (r *tagRepo) GetByName(name string) (*portanizer.Tag, error) {
	var tag portanizer.Tag
	if r.DB.First(&tag, "name = ?", name).RecordNotFound() {
		return nil, portanizer.ErrNotFound
	}
	return &tag, nil
}

func (r *tagRepo) GetList() ([]*portanizer.Tag, error) {
	var tags []*portanizer.Tag
	if err := r.DB.Order("ID ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *tagRepo) Create(tag portanizer.Tag) (*portanizer.Tag, error) {
	if !r.DB.First(&tag, "name = ?", tag.Name).RecordNotFound() {
		return nil, portanizer.ErrExists
	}
	if err := r.DB.Save(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepo) Update(id string, tag portanizer.Tag) (*portanizer.Tag, error) {
	if !r.DB.First(&tag, "name = ?", tag.Name).RecordNotFound() && id != fmt.Sprint(tag.ID) {
		return nil, portanizer.ErrExists
	}
	var updTag portanizer.Tag
	if r.DB.First(&updTag, "id = ?", id).RecordNotFound() {
		return nil, portanizer.ErrNotFound
	}
	if err := r.DB.Model(&updTag).Update(tag).Error; err != nil {
		return nil, err
	}
	return &updTag, nil
}

func (r *tagRepo) Delete(id string) error {
	var tag portanizer.Tag
	if r.DB.First(&tag, "id = ?", id).RecordNotFound() {
		return portanizer.ErrNotFound
	}
	if err := r.DB.Delete(&tag).Error; err != nil {
		return err
	}
	return nil
}
