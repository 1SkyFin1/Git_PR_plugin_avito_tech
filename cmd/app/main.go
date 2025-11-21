package main

import (
	"Git_PR_plugin_avito_tech/cmd/migrate"
	"Git_PR_plugin_avito_tech/config"
	"Git_PR_plugin_avito_tech/internal/app"
	"log"

	"github.com/jmoiron/sqlx"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	log.Println("Starting database migrations...")
	if err := migrate.Run(cfg.PG.URL); err != nil {
		log.Fatalf("Migration error: %v", err)
	}
	log.Println("Database migrations completed successfully")

	log.Println("Connecting to database...")
	db, err := sqlx.Connect("postgres", cfg.PG.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}(db)
	log.Println("Successfully connected to database")

	app.Run(db)
}
