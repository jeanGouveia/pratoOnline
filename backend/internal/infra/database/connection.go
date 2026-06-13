package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	// Para Oracle no futuro, substituir as 2 linhas acima por:
	// "gorm.io/driver/oracle"   (wrapper godror)
)

// Connect é o ÚNICO arquivo a mudar na migração SQLite → Oracle.
// Substitua sqlite.Open("app.db") por oracle.Open(dsn) e pronto.
func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
