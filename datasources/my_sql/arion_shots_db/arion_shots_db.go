package arion_shots_db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
)

const (
	mySQLUsername = "MYSQL_USER"
	mySQLPassword = "MYSQL_PASSWORD"
	mySQLHost     = "MYSQL_HOST"
	mySQLSchema   = "MYSQL_DATABASE"
)

var (
	Client *sql.DB

	username string
	password string
	host     string
	schema   string
)

// init establish the connection with the MySQL database.
func init() {
	var envs map[string]string
	envs, err := godotenv.Read()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username = envs[mySQLUsername]
	password = envs[mySQLPassword]
	host = envs[mySQLHost]
	schema = envs[mySQLSchema]
}

func init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username, password, host, schema)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	Client = db
	log.Println("database successfully configured")
}
