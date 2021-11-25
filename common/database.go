package common

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func InitDb() {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	db, err := sql.Open("mysql", os.Getenv("DB_CONN"))
	if err != nil {
		panic(err)
	}

	query := `CREATE TABLE IF NOT EXISTS search_logs(
		id int primary key auto_increment, 
		created_at datetime default CURRENT_TIMESTAMP,
		transport text, 
		url text
	)`

	_, err = db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating search_logs table", err)
	}

	Db = db
}

func GetDb() *sql.DB {
	if Db != nil {
		return Db
	}

	InitDb()
	return Db
}
