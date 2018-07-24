package main

import (
	"fmt"
	"log"

	app "github.com/dimdiden/portanizer_sop"
	"github.com/dimdiden/portanizer_sop/configure"
	"github.com/dimdiden/portanizer_sop/gorm"
	"github.com/dimdiden/portanizer_sop/http"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const CONFFILE = "./secrets.json"

var (
	tagStore  app.TagStore
	postStore app.PostStore
)

func main() {
	// Create Conf object to use it in starting the server
	c, err := configure.FromFile(CONFFILE)
	if err != nil {
		fmt.Println(err)
		c = configure.Default()
		fmt.Print("Running from the default configuration:\n", c)
	}
	// Open the database
	cs := fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True&loc=Local", c.User, c.DbName)
	db, err := gorm.Open(c.Driver, cs)
	if err != nil {
		log.Fatal("Error opening database:", err)
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
	log.Fatal(http.ListenAndServe(":"+c.Port, server))
}
