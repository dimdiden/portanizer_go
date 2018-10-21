package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dimdiden/portanizer_go/gorm"
	"github.com/dimdiden/portanizer_go/server"
)

func main() {
	// Load the configuration either from environment or from the default values
	c := NewConf()
	fmt.Println("[[> configurator initiated...")
	fmt.Println(c)

	db, err := c.openGormDB()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[[> database connection has been established...")
	defer db.Close()

	// Migrate any changed in structs to DB schema
	gorm.RunMigrations(db)
	// Log each sql query
	db.LogMode(true)
	// Init server, enable logs and run it
	s := server.New(
		c.Secret,
		gorm.NewPostRepo(db),
		gorm.NewTagRepo(db),
		gorm.NewUserRepo(db),
	)
	s.LogEnable()

	fmt.Printf("[[> listening on %v port...", c.APPport)
	log.Fatal(http.ListenAndServe(":"+c.APPport, s))
}
