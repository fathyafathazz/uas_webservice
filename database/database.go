package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
type User struct {
    ID           uint   
    Username     string
    PasswordHash string
    Role         string
}
func InitDB() {
	var err error
	connStr := "root:p3ws@tcp(localhost:3306)/dbwebservice"
	DB, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected")
}