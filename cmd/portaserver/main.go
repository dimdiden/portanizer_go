package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dimdiden/portanizer_go"
	"github.com/dimdiden/portanizer_go/gorm"
	"github.com/dimdiden/portanizer_go/server"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	postRepo portanizer.PostRepo
	tagRepo  portanizer.TagRepo
)

func main() {
	// Load the configuration either from environment or from the default values
	c := NewConf()
	fmt.Print("running configuration:\n", c)

	db, err := c.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("database connection has been established")
	defer db.Close()

	// Migrate any changed in structs to DB schema
	gorm.RunMigrations(db)
	// Log each sql query
	db.LogMode(true)

	// Assigning the store implementation to the server and intiating it
	tagRepo = &gorm.TagRepo{DB: db}
	postRepo = &gorm.PostRepo{DB: db}
	server := server.New(tagRepo, postRepo)
	// Enable the http logs and run
	server.LogEnable()

	log.Fatal(http.ListenAndServe(":"+c.APPport, server))
}
