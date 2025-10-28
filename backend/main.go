package main

import (
	"log"
	"os"

	"shared-expenses-app/db"
	"shared-expenses-app/utils"
)

func main() {
	utils.Loadenv()

	// Open database
	pool, err := db.Connect(utils.Getenv("DB_URL", "postgres://postgres:postgres@localhost:5432/shared_expenses"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Run Migrations
	if err := db.Migrate(pool, os.Getenv("DB_MIGRATIONS_DIR")); err != nil {
		log.Fatal(err)
	}
}
