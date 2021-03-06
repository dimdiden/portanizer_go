package gorm

import (
	"github.com/dimdiden/portanizer_go"
	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"
)

var Open = gorm.Open

type DB = gorm.DB
type Logger = gorm.Logger

func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "initial",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&portanizer.Tag{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable(&portanizer.Tag{}).Error
			},
		},
		{
			ID: "add Post migration",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&portanizer.Post{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable(&portanizer.Post{}).Error
			},
		},
		{
			ID: "add User migration",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&portanizer.User{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable(&portanizer.User{}).Error
			},
		},
		{
			ID: "Add Tokens for User",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&portanizer.User{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable(&portanizer.User{}).Error
			},
		},
	})
	return m.Migrate()
}
