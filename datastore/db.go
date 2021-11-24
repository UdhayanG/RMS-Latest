package datastore

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(host, user, pw, port, dbName string) (*gorm.DB, error) {
	// check env variables
	//environment.CheckEnvironmentVariables(requiredEnvironmentVariablesForMySQL)
	// env variables

	//sda
	//fmt.Printf("============:%s %s %s %s %s", host, user, pw, port, dbName)
	// build connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pw, host, port, dbName)
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "root", "localhost", "3306", "CentralUserNew")
	// connect to db
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})

}

/*func NewDB() (*gorm.DB, error) {
	DBMS := "mysql"
	mySqlConfig := &mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "rms-golang",
		AllowNativePasswords: true,
		Params: map[string]string{
			"parseTime": "true",
		},
	}

	return gorm.Open(DBMS, mySqlConfig.FormatDSN())
}*/
