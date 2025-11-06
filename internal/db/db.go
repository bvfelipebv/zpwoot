package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver

	"zpwoot/internal/config"
	"zpwoot/pkg/logger"
)

var DB *sql.DB

func InitDB() error {
	dsn := config.GetDatabaseDSN()
	driver := config.AppConfig.DatabaseDriver



	// Abrir conexão
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Testar conexão
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	DB = db

	// Executar migrações automaticamente
	if err := RunMigrations(context.Background()); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
