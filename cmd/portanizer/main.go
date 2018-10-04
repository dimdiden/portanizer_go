package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dimdiden/portanizer_go"
	"github.com/dimdiden/portanizer_go/gorm"
	"github.com/dimdiden/portanizer_go/server"
)

var (
	postRepo portanizer.PostRepo
	tagRepo  portanizer.TagRepo
)

func main() {
	// Load the configuration either from environment or from the default values
	c := NewConf()
	fmt.Println("[[> configurator initiated...")
	fmt.Print(c)

	db, err := c.OpenGormDB()
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
	s := server.New(gorm.NewPostRepo(db), gorm.NewTagRepo(db))
	s.LogEnable()
	fmt.Printf("[[> listening on %v port...", c.APPport)
	log.Fatal(http.ListenAndServe(":"+c.APPport, s))
}
