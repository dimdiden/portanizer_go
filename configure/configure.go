package configure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	HOST = ""
	PORT = "8080"

	DRIVER = "mysql"
	DBNAME = "portanizer_sop"
	USER   = "root"
	PSWD   = ""
)

type Conf struct {
	// Host and port to run the server with
	Host string
	Port string
	// The information for DB connection
	Driver string
	DbName string
	User   string
	Pswd   string
}

func (c Conf) String() string {
	s := fmt.Sprintf("HOST: %v\nPORT: %v\nDRIVER: %v\nDBNAME: %v\nUSER: %v\nPSWD: %v\n",
		c.Host, c.Port, c.Driver, c.DbName, c.User, c.Pswd)
	return s
}

func FromFile(file string) (*Conf, error) {
	var conf Conf
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(dat, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func Default() *Conf {
	return &Conf{
		Host:   HOST,
		Port:   PORT,
		Driver: DRIVER,
		DbName: DBNAME,
		User:   USER,
		Pswd:   PSWD,
	}
}
