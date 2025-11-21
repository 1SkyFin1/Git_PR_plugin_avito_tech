package migrate

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Run(databaseURL string) error {
	m, err := migrate.New("file://./resources/migrations", databaseURL)
	if err != nil {
		return err
	}
	defer func() {
		sourceErr, databaseErr := m.Close()
		if sourceErr != nil {
			log.Printf("Error closing migration source: %v", sourceErr)
		}
		if databaseErr != nil {
			log.Printf("Error closing migration database: %v", databaseErr)
		}
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		log.Printf("Warning: could not get migration version: %v", err)
	} else if err == nil {
		status := "clean"
		if dirty {
			status = "dirty"
		}
		log.Printf("Current migration version: %d (status: %s)", version, status)
	} else {
		log.Println("No migrations have been applied yet")
	}

	return nil
}
