package db

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"strings"

	"zpwoot/pkg/logger"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migration representa uma migração de banco de dados
type Migration struct {
	Version int
	Name    string
	UpSQL   string
	DownSQL string
}

// RunMigrations executa todas as migrações pendentes
func RunMigrations(ctx context.Context) error {
	logger.Log.Info().Msg("Starting database migrations...")

	// Criar tabela de controle de migrações
	if err := createMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Carregar migrações dos arquivos
	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Obter migrações já aplicadas
	appliedVersions, err := getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Executar migrações pendentes
	for _, migration := range migrations {
		if _, applied := appliedVersions[migration.Version]; applied {
			logger.Log.Debug().
				Int("version", migration.Version).
				Str("name", migration.Name).
				Msg("Migration already applied, skipping")
			continue
		}

		logger.Log.Info().
			Int("version", migration.Version).
			Str("name", migration.Name).
			Msg("Applying migration")

		if err := applyMigration(ctx, migration); err != nil {
			return fmt.Errorf("failed to apply migration %d (%s): %w", migration.Version, migration.Name, err)
		}

		logger.Log.Info().
			Int("version", migration.Version).
			Str("name", migration.Name).
			Msg("Migration applied successfully")
	}

	logger.Log.Info().Msg("All migrations completed successfully")
	return nil
}

// createMigrationsTable cria a tabela de controle de migrações
func createMigrationsTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`

	_, err := DB.ExecContext(ctx, query)
	return err
}

// loadMigrations carrega todas as migrações dos arquivos embedded
func loadMigrations() ([]Migration, error) {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	migrationsMap := make(map[int]*Migration)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !strings.HasSuffix(filename, ".sql") {
			continue
		}

		// Parse filename: 001_create_sessions.up.sql ou 001_create_sessions.down.sql
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			continue
		}

		var version int
		fmt.Sscanf(parts[0], "%d", &version)

		isUp := strings.HasSuffix(filename, ".up.sql")
		isDown := strings.HasSuffix(filename, ".down.sql")

		if !isUp && !isDown {
			continue
		}

		// Ler conteúdo do arquivo
		content, err := migrationsFS.ReadFile("migrations/" + filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Extrair nome da migração
		name := strings.TrimSuffix(strings.TrimSuffix(filename, ".up.sql"), ".down.sql")
		name = strings.Join(parts[1:], "_")
		name = strings.TrimSuffix(name, ".up")
		name = strings.TrimSuffix(name, ".down")

		if migrationsMap[version] == nil {
			migrationsMap[version] = &Migration{
				Version: version,
				Name:    name,
			}
		}

		if isUp {
			migrationsMap[version].UpSQL = string(content)
		} else {
			migrationsMap[version].DownSQL = string(content)
		}
	}

	// Converter map para slice e ordenar
	migrations := make([]Migration, 0, len(migrationsMap))
	for _, m := range migrationsMap {
		if m.UpSQL != "" { // Só adicionar se tiver SQL de UP
			migrations = append(migrations, *m)
		}
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// getAppliedMigrations retorna as versões de migrações já aplicadas
func getAppliedMigrations(ctx context.Context) (map[int]bool, error) {
	query := `SELECT version FROM schema_migrations ORDER BY version`

	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// applyMigration aplica uma migração específica
func applyMigration(ctx context.Context, migration Migration) error {
	// Iniciar transação
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Executar SQL da migração
	if _, err := tx.ExecContext(ctx, migration.UpSQL); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// Registrar migração como aplicada
	insertQuery := `
		INSERT INTO schema_migrations (version, name, applied_at)
		VALUES ($1, $2, NOW())
	`
	if _, err := tx.ExecContext(ctx, insertQuery, migration.Version, migration.Name); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit da transação
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	return nil
}

// RollbackMigration reverte uma migração específica
func RollbackMigration(ctx context.Context, version int) error {
	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	var migration *Migration
	for _, m := range migrations {
		if m.Version == version {
			migration = &m
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("migration version %d not found", version)
	}

	if migration.DownSQL == "" {
		return fmt.Errorf("migration version %d has no down SQL", version)
	}

	logger.Log.Info().
		Int("version", version).
		Str("name", migration.Name).
		Msg("Rolling back migration")

	// Iniciar transação
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Executar SQL de rollback
	if _, err := tx.ExecContext(ctx, migration.DownSQL); err != nil {
		return fmt.Errorf("failed to execute rollback SQL: %w", err)
	}

	// Remover registro da migração
	deleteQuery := `DELETE FROM schema_migrations WHERE version = $1`
	if _, err := tx.ExecContext(ctx, deleteQuery, version); err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	// Commit da transação
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback: %w", err)
	}

	logger.Log.Info().
		Int("version", version).
		Str("name", migration.Name).
		Msg("Migration rolled back successfully")

	return nil
}
