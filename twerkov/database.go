package twerkov

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

// Database contains information about the MySQL database used for the application
type Database struct {
	Handle *sql.DB
}

// CreateDatabaseConnection estabilishes a connection to a MySQL database
func CreateDatabaseConnection(config Config) (db *Database, err error) {
	db = &Database{}

	dbConfig := mysql.NewConfig()
	dbConfig.Net = "tcp"
	dbConfig.Addr = config.MySQLHostname
	dbConfig.DBName = config.MySQLDatabase
	dbConfig.User = config.MySQLUsername
	dbConfig.Passwd = config.MySQLPassword

	db.Handle, err = sql.Open("mysql", dbConfig.FormatDSN())
	return
}
