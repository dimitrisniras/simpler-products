package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func Init(log *logrus.Logger, user, password, host, port, dbName string) (*sql.DB, error) {
	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)
	db, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		log.Errorf("Error connecting to database: %v", err)
		return nil, err
	}

	// check if Database connection is established
	err = db.Ping()
	if err != nil {
		log.Errorf("Error checking database connection: %v", err)
		return nil, err
	}

	return db, nil
}
