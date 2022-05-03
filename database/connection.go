package database

import (
	"database/sql"
	"fmt"
	"log"
	// "os"

	"github.com/joho/godotenv"
	_"github.com/lib/pq"
)

func CreateConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := `postgres://admin:Xd8UCzVl5Z2U3C6IOJwONKZgRVWzTqz@34.142.140.9:5432/dev-wallet?sslmode=disable`
	// db, err := sql.Open("postgres", os.ExpandEnv("$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSLMODE"))
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database")

	return db
}