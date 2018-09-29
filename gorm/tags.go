package gorm

import (
	"fmt"

	"github.com/dimdiden/portanizer_go"
	"github.com/jinzhu/gorm"
)

type TagRepo struct {
	DB *gorm.DB
}

func (r *TagRepo) GetByID(id string) (*portanizer.Tag, error) {
	var tag portanizer.Tag
	if r.DB.First(&tag, "id = ?", id).RecordNotFound() {
		return nil, portanizer.ErrNotFound
	}
	return &tag, nil
}

func (r *TagRepo) GetByName(name string) (*portanizer.Tag, error) {
	var tag portanizer.Tag
	if r.DB.First(&tag, "name = ?", name).RecordNotFound() {
		return nil, portanizer.ErrNotFound
	}
	return &tag, nil
}

func (r *TagRepo) GetList() ([]*portanizer.Tag, error) {
	var tags []*portanizer.Tag
	if err := r.DB.Order("ID ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *TagRepo) Create(tag portanizer.Tag) (*portanizer.Tag, error) {
	if !r.DB.First(&tag, "name = ?", tag.Name).RecordNotFound() {
		return nil, portanizer.ErrExists
	}
	if err := r.DB.Save(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepo) Update(id string, tag portanizer.Tag) (*portanizer.Tag, error) {
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

func (r *TagRepo) Delete(id string) error {
	var tag portanizer.Tag
	if r.DB.First(&tag, "id = ?", id).RecordNotFound() {
		return portanizer.ErrNotFound
	}
	if err := r.DB.Delete(&tag).Error; err != nil {
		return err
	}
	return nil
}
