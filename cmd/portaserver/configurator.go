package main

import (
	"fmt"
	"os"

	"github.com/dimdiden/portanizer_sop/gorm"
)

const (
	APP_PORT = "8080"

	DB_HOST   = "127.0.0.1"
	DB_DRIVER = "mysql"
	DB_NAME   = "portanizer"
	DB_USER   = "root"
	DB_PSWD   = ""
)

var CONFLIST = map[string]string{
	"APP_PORT":  APP_PORT,
	"DB_HOST":   DB_HOST,
	"DB_DRIVER": DB_DRIVER,
	"DB_NAME":   DB_NAME,
	"DB_USER":   DB_USER,
	"DB_PSWD":   DB_PSWD,
}

type Conf struct {
	// Host and port to run the server with
	APPport string
	// The information for DB connection
	DBhost   string
	DBdriver string
	DBname   string
	DBuser   string
	DBpswd   string
}

func (c Conf) String() string {
	s := fmt.Sprintf("APP_PORT: %v\nDB_HOST: %v\nDB_DRIVER: %v\nDB_NAME: %v\nDB_USER: %v\nDB_PSWD: %v\n",
		c.APPport, c.DBhost, c.DBdriver, c.DBname, c.DBuser, c.DBpswd)
	return s
}

func Default() *Conf {
	return &Conf{
		APPport:  CONFLIST["APP_PORT"],
		DBhost:   CONFLIST["DB_HOST"],
		DBdriver: CONFLIST["DB_DRIVER"],
		DBname:   CONFLIST["DB_NAME"],
		DBuser:   CONFLIST["DB_USER"],
		DBpswd:   CONFLIST["DB_PSWD"],
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
	}
}

func getOpt(opt string) string {
	val, ok := os.LookupEnv(opt)
	if !ok {
		return CONFLIST[opt]
	}
	return val
}

func (c *Conf) OpenDB() (*gorm.DB, error) {
	cs := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", c.DBuser, c.DBpswd, c.DBhost, c.DBname)
	db, err := gorm.Open(c.DBdriver, cs)
	if err != nil {
		return nil, fmt.Errorf("unable to open app db: %v", err)

	}
	return db, nil
}
