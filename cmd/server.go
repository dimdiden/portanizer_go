package main

import (
	"fmt"
	"log"

	app "github.com/dimdiden/portanizer_sop"
	"github.com/dimdiden/portanizer_sop/gorm"
	"github.com/dimdiden/portanizer_sop/http"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	driver = "mysql"
	host   = "localhost"
	// port   = "8080"
	user = "root"
	// password = "your-password"
	dbname = "portanizer_sop"
)

var (
	tagStore app.TagStore
)

func main() {
	cs := fmt.Sprintf("%s:@/%s?charset=utf8&parseTime=True&loc=Local", user, dbname)
	db, err := gorm.Open(driver, cs)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()
	gorm.RunMigrations(db)

	tagStore = &gorm.TagService{DB: db}

	server := http.NewServer(tagStore)
	log.Fatal(http.ListenAndServe(":8080", server))
}
