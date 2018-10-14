package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/dimdiden/portanizer_go/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	// Default values
	APP_PORT = "8080"

	DB_HOST   = "127.0.0.1"
	DB_DRIVER = "mysql"
	DB_NAME   = "portanizer"
	DB_USER   = "root"
	DB_PSWD   = ""

	IS_DEBUG = "OFF"
)

type Conf struct {
	// Host and port to run the server with
	APPport string
	// The information for DB connection
	DBhost   string
	DBdriver string
	DBname   string
	DBuser   string
	DBpswd   string

	IsDebug string
}

var conflist = map[string]string{
	"APP_PORT":  APP_PORT,
	"DB_HOST":   DB_HOST,
	"DB_DRIVER": DB_DRIVER,
	"DB_NAME":   DB_NAME,
	"DB_USER":   DB_USER,
	"DB_PSWD":   DB_PSWD,
	"IS_DEBUG":  IS_DEBUG,
}

func (c Conf) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "[[> running configuration")
	fmt.Fprint(w, "    IS_DEBUG:\t"+c.IsDebug)
	if c.IsDebug == "ON" {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "    APP_PORT:\t%v\n    DB_HOST:\t%v\n    DB_DRIVER:\t%v\n    DB_NAME:\t%v\n    DB_USER:\t%v\n    DB_PSWD:\t%v", // without \n because of w - specific writer
			c.APPport, c.DBhost, c.DBdriver, c.DBname, c.DBuser, c.DBpswd)
	}
	w.Flush()
	return buf.String()
}

func NewDefaultConf() *Conf {
	return &Conf{
		APPport:  conflist["APP_PORT"],
		DBhost:   conflist["DB_HOST"],
		DBdriver: conflist["DB_DRIVER"],
		DBname:   conflist["DB_NAME"],
		DBuser:   conflist["DB_USER"],
		DBpswd:   conflist["DB_PSWD"],
		IsDebug:  conflist["IS_DEBUG"],
	}
}

func NewConf() *Conf {
	return &Conf{
		APPport:  getOpt("APP_PORT"),
		DBhost:   getOpt("DB_HOST"),
		DBdriver: getOpt("DB_DRIVER"),
		DBname:   getOpt("DB_NAME"),
		DBuser:   getOpt("DB_USER"),
		DBpswd:   getOpt("DB_PSWD"),
		IsDebug:  getOpt("IS_DEBUG"),
	}
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

	switch c.DBdriver {
	case "mysql":
		cparams = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", c.DBuser, c.DBpswd, c.DBhost, c.DBname)
	case "sqlite3":
		cparams = "./sqlite.db"
	default:
		return nil, errors.New("unsupported dialect for database")
	}

	db, err := gorm.Open(c.DBdriver, cparams)
	if err != nil {
		return nil, fmt.Errorf("unable to open app db: %v", err)
	}
	return db, nil
}
