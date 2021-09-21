package managementapi

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

var db *sqlx.DB

func init() {
	dbDriver := os.Getenv("DATABASE_DRIVER")
	dbSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_SSLMODE"))

	var err error
	if db, err = sqlx.Open(dbDriver, dbSource); err != nil {
		log.Fatalf("db connection error: %v", err)
	}
}
