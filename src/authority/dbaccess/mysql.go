package dbaccess

import (
//	"database/sql"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"common/logging"
)

var db *sqlx.DB

func DB() *sqlx.DB {
	return db
}

func InitMysql(mysqlurl string) {
	if ok, err := regexp.MatchString("^mysql://.*:.*@.*/.*$", mysqlurl); ok == false || err != nil {
		logging.Error("mysql config syntax err:mysql_zone,%s,shutdown", mysqlurl)
		panic("InitMysql conf error")
		//return nil
	}
	mysqlurl = strings.Replace(mysqlurl, "mysql://", "", 1)

	var err error
	db, err = sqlx.Open("mysql", mysqlurl)
	if err != nil {
		logging.Error("InitMysql failed mysqlurl=" + mysqlurl + ",err=" + err.Error())
		panic("InitMysql failed mysqlurl=" + mysqlurl)
		//return nil
	} else {
		logging.Info("mysql conn ok:%s", mysqlurl)
	}
	//return db
}
