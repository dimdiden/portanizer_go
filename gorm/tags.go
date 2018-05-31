package gorm

import (
	app "github.com/dimdiden/portanizer_sop"
	"github.com/jinzhu/gorm"
)

type TagService struct {
	DB *gorm.DB
}

func (s *TagService) GetByID(id string) (*app.Tag, error) {
	var tag app.Tag
	if s.DB.First(&tag, "id = ?", id).RecordNotFound() {
		return nil, app.ErrNotFound
	}
	return &tag, nil
}

func (s *TagService) GetByName(name string) (*app.Tag, error) {
	var tag app.Tag
	if s.DB.First(&tag, "name = ?", name).RecordNotFound() {
		return nil, app.ErrNotFound
	}
	return &tag, nil
}

func (s *TagService) GetList() ([]*app.Tag, error) {
	var tags []*app.Tag
	s.DB.Order("ID ASC").Find(&tags)
	return tags, nil
}

func (s *TagService) Create(tag app.Tag) (*app.Tag, error) {
	if err := s.DB.Save(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *TagService) Update(id string, tag app.Tag) (*app.Tag, error) {
	if s.DB.First(&tag, "id = ?", id).RecordNotFound() {
		return nil, app.ErrNotFound
	}
	if err := s.DB.Model(&tag).Updates(app.Tag{Name: tag.Name}).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *TagService) Delete(id string) error {
	var tag app.Tag
	if s.DB.First(&tag, "id = ?", id).RecordNotFound() {
		return app.ErrNotFound
	}
	if err := s.DB.Delete(&tag).Error; err != nil {
		return err
	}
	return nil
}
