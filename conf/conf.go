package conf

import (
	"fmt"
	"os"
)

const (
	APP_PORT = "8080"

	DB_HOST   = "127.0.0.1"
	DB_DRIVER = "mysql"
	DB_NAME   = "portanizer"
	DB_USER   = "root"
	DB_PSWD   = "Ltleirf180586"
)

var CONFLIST = map[string]string{
	"APPport":  APP_PORT,
	"DBhost":   DB_HOST,
	"DBdriver": DB_DRIVER,
	"DBname":   DB_NAME,
	"DBuser":   DB_USER,
	"DBpswd":   DB_PSWD,
}

type Conf struct {
	// Host and port to run the server with
	APPport string
	DBhost  string
	// The information for DB connection
	DBdriver string
	DBname   string
	DBuser   string
	DBpswd   string
}

func (c Conf) String() string {
	s := fmt.Sprintf("APP_PORT: %v\nDBhost: %v\nDBdriver: %v\nDBname: %v\nDBuser: %v\nDBpswd: %v\n",
		c.APPport, c.DBhost, c.DBdriver, c.DBname, c.DBuser, c.DBpswd)
	return s
}

func Default() *Conf {
	return &Conf{
		APPport:  CONFLIST["APPport"],
		DBhost:   CONFLIST["DBhost"],
		DBdriver: CONFLIST["DBdriver"],
		DBname:   CONFLIST["DBname"],
		DBuser:   CONFLIST["DBuser"],
		DBpswd:   CONFLIST["DBpswd"],
	}
}

func Get() *Conf {
	conf := &Conf{
		APPport:  getOpt("APPport"),
		DBhost:   getOpt("DBhost"),
		DBdriver: getOpt("DBdriver"),
		DBname:   getOpt("DBname"),
		DBuser:   getOpt("DBuser"),
		DBpswd:   getOpt("DBpswd"),
	}
	return conf
}

func getOpt(opt string) string {
	val, ok := os.LookupEnv(opt)
	if !ok {
		return CONFLIST[opt]
	}
	return val
}
