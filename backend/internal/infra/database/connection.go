package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	DSN string // SQLite: "app.db" | Oracle futuro: "user/pass@host:1521/svc"
}

func Connect(cfg DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Warn),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("database.Connect: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("database.Connect: obter sql.DB: %w", err)
	}

	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA foreign_keys=ON;",
		"PRAGMA cache_size=-65536;",
		"PRAGMA synchronous=NORMAL;",
	}
	for _, p := range pragmas {
		if _, err = sqlDB.Exec(p); err != nil {
			return nil, fmt.Errorf("database.Connect: pragma '%s': %w", p, err)
		}
	}

	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	return db, nil
}
