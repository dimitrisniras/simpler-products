package database

import (
	"database/sql"
	"fmt"
	"time"

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

	// Start a goroutine to periodically check the connection
	go checkDBConnection(db, dbConnectionString, log)

	return db, nil
}

func checkDBConnection(db *sql.DB, dbConnectionString string, log *logrus.Logger) {
	for {
		// Sleep for 5s
		time.Sleep(5 * time.Second)

		// Check the connection
		err := db.Ping()

		if err != nil {
			log.Errorf("Database connection lost: %v", err)

			// Attempt to reconnect
			for {
				log.Info("Attempting to reconnect to the database...")
				db, err = sql.Open("mysql", dbConnectionString)
				if err != nil {
					log.Errorf("Reconnection failed: %v", err)
					time.Sleep(5 * time.Second)
					continue
				}

				err = db.Ping()
				if err != nil {
					log.Errorf("Reconnection ping failed: %v", err)
					time.Sleep(5 * time.Second)
					continue
				}

				log.Info("Successfully reconnected to the database!")
				break
			}
		}
	}
}
