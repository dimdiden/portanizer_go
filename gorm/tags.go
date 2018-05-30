package gorm

import (
	"errors"

	app "github.com/dimdiden/portanizer_sop"
	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"
)

type TagService struct {
	DB *gorm.DB
}

var Open = gorm.Open

func (s *TagService) GetByID(id string) (*app.Tag, error) {
	var tag app.Tag
	if s.DB.First(&tag, "id = ?", id).RecordNotFound() {
		return nil, errors.New("Record not found")
	}
	return &tag, nil
}

func (s *TagService) GetByName(name string) (*app.Tag, error) {
	var tag app.Tag
	if s.DB.First(&tag, "name = ?", name).RecordNotFound() {
		return nil, errors.New("Record not found")
	}
	return &tag, nil
}

func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "initial",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&app.Tag{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable(&app.Tag{}).Error
			},
		},
	})
	return m.Migrate()
}
