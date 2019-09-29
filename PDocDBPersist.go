package PttUtils

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)
import _ "github.com/go-sql-driver/mysql"

var g_DbObject *DBObject

//Fetch from BossBot, can it be merged to personal module?
type DBObject struct {
	Host       string
	User       string
	Passwd     string
	connString string
	conn       *sql.DB
}


func NewDBObject(host string, user string, passwd string) (*DBObject, error) {
	ret := DBObject{
		Host:   host,
		User:   user,
		Passwd: passwd,
	}

	ret.connString = fmt.Sprintf("%s:%s@tcp(%s)/apps?charset=utf8&loc=Asia%%2FTaipei&parseTime=true", user, passwd, host)

	db, err := sql.Open("mysql", ret.connString)
	log.Debugf("Attempting logging in with Server String : %s", ret.connString)
	if err != nil {
		wrapped := errors.Wrap(err, "Error while initialization with sql string : "+ret.connString)
		return nil, wrapped
	}
	log.Debugf("Attempting pinging : %s", ret.connString)
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "Fail to ping server : "+ret.connString)
	}

	ret.conn = db

	return &ret, nil
}

func (db *DBObject) GetDB() *sql.DB {
	return db.conn
}

func InitDBObject(host string, user string, passwd string) {
	g_DbObject, _ = NewDBObject(host, user, passwd)
}

func GetDBObject() *DBObject{
	return g_DbObject
}