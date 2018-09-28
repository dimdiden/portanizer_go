package main

import (
	"fmt"
	"log"
	"time"

	app "github.com/dimdiden/portanizer_sop"
	"github.com/dimdiden/portanizer_sop/gorm"
	"github.com/dimdiden/portanizer_sop/http"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	tagStore  app.TagStore
	postStore app.PostStore
)

func main() {
	// Load the configuration either from environment or from the default values
	c := NewConf()
	fmt.Print("Running configuration:\n", c)

	// Open the GORM istance of the database
	var db *gorm.DB
	var err error
	for {
		db, err = c.OpenDB()
		if err == nil {
			break
		}
		log.Println("Error opening database:", err)
		time.Sleep(1 * time.Second)
	}

	defer db.Close()
	// Migrate any changed in structs to DB schema
	gorm.RunMigrations(db)
	// Log each sql query
	db.LogMode(true)

	// Assigning the store implementation to the server and intiating it
	tagStore = &gorm.TagService{DB: db}
	postStore = &gorm.PostService{DB: db}
	server := http.NewServer(tagStore, postStore)
	// Enable the http logs and run
	server.LogHttpEnable()
	log.Fatal(http.ListenAndServe(":"+c.APPport, server))
}
