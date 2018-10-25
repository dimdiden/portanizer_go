package main

import (
	"log"

	"github.com/dimdiden/portanizer_go/gorm"
)

func main() {
	// Load the configuration either from environment or from the default values
	c := newConf()
	// Open the database
	db, err := c.openGormDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	// Migrate any changed in structs to DB schema
	gorm.RunMigrations(db)
	// Init server and run it
	s := c.openGormServer(db)
	log.Fatal(s.Run(c.APPport))
}
