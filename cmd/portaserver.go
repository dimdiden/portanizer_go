package main

import (
	"fmt"
	"log"

	app "github.com/dimdiden/portanizer_sop"
	conf "github.com/dimdiden/portanizer_sop/conf"
	"github.com/dimdiden/portanizer_sop/gorm"
	"github.com/dimdiden/portanizer_sop/http"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// const CONFFILE = "./secrets.json"

var (
	tagStore  app.TagStore
	postStore app.PostStore
)

func main() {
	// Load the configuration either from environment or from the default values
	c := conf.Get()
	fmt.Print("Running configuration:\n", c)

	// Open the GORM istance of the database
	cs := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", c.DBuser, c.DBpswd, c.DBhost, c.DBname)
	db, err := gorm.Open(c.DBdriver, cs)
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
	log.Fatal(http.ListenAndServe(":"+c.APPport, server))
}
