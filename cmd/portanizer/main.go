package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dimdiden/portanizer_go"
	"github.com/dimdiden/portanizer_go/gorm"
	"github.com/dimdiden/portanizer_go/server"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	// Init server, enable logs and run it
	s := server.New(gorm.NewPostRepo(db), gorm.NewTagRepo(db))
	s.LogEnable()
	log.Fatal(http.ListenAndServe(":"+c.APPport, s))
}
