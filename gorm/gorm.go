package gorm

import (
	app "github.com/dimdiden/portanizer_sop"
	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"
)

var Open = gorm.Open

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
		{
			ID: "add Post migration",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&app.Post{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable(&app.Post{}).Error
			},
		},
	})
	return m.Migrate()
}
