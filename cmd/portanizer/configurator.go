package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"text/tabwriter"

	"github.com/dimdiden/portanizer_go/gorm"
	"github.com/dimdiden/portanizer_go/server"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	APP_PORT = "8080"

	DB_HOST   = "127.0.0.1"
	DB_DRIVER = "sqlite3"
	DB_NAME   = "portanizer"
	DB_USER   = "root"
	DB_PSWD   = ""

	DEBUG   = "OFF"
	ASECRET = "ACCESS_SECRET_KEY"
	RSECRET = "REFRESH_SECRET_KEY"
)

var conflist = map[string]string{
	"APP_PORT":  APP_PORT,
	"DB_HOST":   DB_HOST,
	"DB_DRIVER": DB_DRIVER,
	"DB_NAME":   DB_NAME,
	"DB_USER":   DB_USER,
	"DB_PSWD":   DB_PSWD,
	"DEBUG":     DEBUG,
	"ASECRET":   ASECRET,
	"RSECRET":   RSECRET,
}

type Conf struct {
	Debug string
	// Host and port to run the server with
	APPport string
	// The information for DB connection
	DBhost   string
	DBdriver string
	DBname   string
	DBuser   string
	DBpswd   string

	ASecret string
	RSecret string

	logout io.Writer
}

func (c Conf) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "[[> running configuration")
	fmt.Fprint(w, "    DEBUG:\t"+c.Debug)
	if c.Debug == "ON" {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "    ASECRET:\t%v\n    RSECRET:\t%v\n    APP_PORT:\t%v\n    DB_HOST:\t%v\n    DB_DRIVER:\t%v\n    DB_NAME:\t%v\n    DB_USER:\t%v\n    DB_PSWD:\t%v", // without \n because of w - specific writer
			c.ASecret, c.RSecret, c.APPport, c.DBhost, c.DBdriver, c.DBname, c.DBuser, c.DBpswd)
	}
	w.Flush()
	return buf.String()
}

func newConf() *Conf {
	conf := &Conf{
		APPport:  getOpt("APP_PORT"),
		DBhost:   getOpt("DB_HOST"),
		DBdriver: getOpt("DB_DRIVER"),
		DBname:   getOpt("DB_NAME"),
		DBuser:   getOpt("DB_USER"),
		DBpswd:   getOpt("DB_PSWD"),
		Debug:    getOpt("DEBUG"),
		ASecret:  getOpt("ASECRET"),
		RSecret:  getOpt("RSECRET"),
		logout:   os.Stdout,
	}
	fmt.Fprintln(conf.logout, "[[> configurator initiated...")
	fmt.Fprintln(conf.logout, conf)
	return conf
}

func getOpt(opt string) string {
	val, ok := os.LookupEnv(opt)
	if !ok {
		// return default value
		return conflist[opt]
	}
	return val
}

func (c *Conf) openGormDB() (*gorm.DB, error) {
	var cparams string
	// get the connection string
	switch c.DBdriver {
	case "mysql":
		cparams = fmt.Sprintf(
			"%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local",
			c.DBuser, c.DBpswd, c.DBhost, c.DBname)
	case "sqlite3":
		cparams = "./sqlite.db"
	default:
		return nil, errors.New("unsupported dialect for database")
	}
	// open db and set logger to use conf.logout
	db, err := gorm.Open(c.DBdriver, cparams)
	if err != nil {
		return nil, fmt.Errorf("unable to open app db: %v", err)
	}
	logger := gorm.Logger{log.New(c.logout, "\r\n", 0)}
	db.SetLogger(logger)
	if c.Debug == "ON" {
		db.LogMode(true)
	}
	fmt.Fprintln(c.logout, "[[> database connection has been established...")
	return db, nil
}

func (c *Conf) openGormServer(db *gorm.DB) *server.Server {
	server.ASecret = []byte(c.ASecret)
	server.RSecret = []byte(c.RSecret)
	s := server.New(
		c.logout,
		gorm.NewPostRepo(db),
		gorm.NewTagRepo(db),
		gorm.NewUserRepo(db),
	)
	return s
}
