package migration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/yeferson59/db-migration-cli/internal/config"
	"github.com/yeferson59/db-migration-cli/internal/database"
)

// Service handles migration operations
type Service struct {
	cfg *config.Config
	db  *database.DB
}

// Migration represents a database migration
type Migration struct {
	Version   string
	Name      string
	UpSQL     string
	DownSQL   string
	Applied   bool
	AppliedAt *time.Time
}

// NewService creates a new migration service
func NewService() *Service {
	return &Service{}
}

// Init initializes the migration environment
func (s *Service) Init() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	s.cfg = cfg

	// Create migrations directory if it doesn't exist
	if err := os.MkdirAll(s.cfg.Migrations.Dir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Connect to database
	db, err := database.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Create schema_migrations table
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`

	if err := db.Exec(context.Background(), createTableSQL); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	return nil
}

// Create creates a new migration file
func (s *Service) Create(name string) (string, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}
	s.cfg = cfg

	// Generate timestamp-based version
	version := time.Now().Format("20060102150405")

	// Create up migration file
	upFilename := fmt.Sprintf("%s_%s.up.sql", version, name)
	upPath := filepath.Join(s.cfg.Migrations.Dir, upFilename)

	upContent := fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n\n-- Write your UP migration here\n\n",
		name, time.Now().Format(time.RFC3339))

	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create up migration file: %w", err)
	}

	// Create down migration file
	downFilename := fmt.Sprintf("%s_%s.down.sql", version, name)
	downPath := filepath.Join(s.cfg.Migrations.Dir, downFilename)

	downContent := fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n\n-- Write your DOWN migration here\n\n",
		name, time.Now().Format(time.RFC3339))

	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		return "", fmt.Errorf("failed to create down migration file: %w", err)
	}

	return fmt.Sprintf("%s_%s", version, name), nil
}

// Up applies all pending migrations
func (s *Service) Up() (int, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return 0, fmt.Errorf("failed to load config: %w", err)
	}
	s.cfg = cfg

	// Connect to database
	db, err := database.New(cfg)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()
	s.db = db

	// Get all migrations
	migrations, err := s.getMigrations()
	if err != nil {
		return 0, err
	}

	// Apply pending migrations
	count := 0
	for _, m := range migrations {
		if !m.Applied {
			if err := s.applyMigration(m); err != nil {
				return count, fmt.Errorf("failed to apply migration %s: %w", m.Version, err)
			}
			count++
			fmt.Printf("  Applied: %s_%s\n", m.Version, m.Name)
		}
	}

	return count, nil
}

// Down rolls back the last migration
func (s *Service) Down() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	s.cfg = cfg

	// Connect to database
	db, err := database.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()
	s.db = db

	// Get all migrations
	migrations, err := s.getMigrations()
	if err != nil {
		return err
	}

	// Find the last applied migration
	var lastMigration *Migration
	for i := len(migrations) - 1; i >= 0; i-- {
		if migrations[i].Applied {
			lastMigration = &migrations[i]
			break
		}
	}

	if lastMigration == nil {
		return fmt.Errorf("no migrations to rollback")
	}

	// Rollback the migration
	if err := s.rollbackMigration(*lastMigration); err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", lastMigration.Version, err)
	}

	fmt.Printf("  Rolled back: %s_%s\n", lastMigration.Version, lastMigration.Name)
	return nil
}

// Status returns the status of all migrations
func (s *Service) Status() ([]string, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	s.cfg = cfg

	// Connect to database
	db, err := database.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()
	s.db = db

	// Get all migrations
	migrations, err := s.getMigrations()
	if err != nil {
		return nil, err
	}

	var status []string
	for _, m := range migrations {
		statusStr := "[ ]"
		if m.Applied {
			statusStr = "[✓]"
		}
		status = append(status, fmt.Sprintf("%s %s_%s", statusStr, m.Version, m.Name))
	}

	if len(status) == 0 {
		status = append(status, "No migrations found")
	}

	return status, nil
}

// getMigrations returns all migrations sorted by version
func (s *Service) getMigrations() ([]Migration, error) {
	// Read migration files
	files, err := os.ReadDir(s.cfg.Migrations.Dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	migrationMap := make(map[string]*Migration)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		version, migrationName, isUp, err := parseMigrationFilename(name)
		if err != nil {
			continue // Skip invalid files
		}

		if _, exists := migrationMap[version]; !exists {
			migrationMap[version] = &Migration{
				Version: version,
				Name:    migrationName,
			}
		}

		content, err := os.ReadFile(filepath.Join(s.cfg.Migrations.Dir, name))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", name, err)
		}

		if isUp {
			migrationMap[version].UpSQL = string(content)
		} else {
			migrationMap[version].DownSQL = string(content)
		}
	}

	// Get applied migrations from database
	appliedVersions, err := s.getAppliedMigrations()
	if err != nil {
		return nil, err
	}

	// Convert map to slice and mark applied migrations
	var migrations []Migration
	for _, m := range migrationMap {
		if appliedAt, applied := appliedVersions[m.Version]; applied {
			m.Applied = true
			m.AppliedAt = &appliedAt
		}
		migrations = append(migrations, *m)
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// getAppliedMigrations returns a map of applied migration versions
func (s *Service) getAppliedMigrations() (map[string]time.Time, error) {
	rows, err := s.db.Query(
		context.Background(),
		"SELECT version, applied_at FROM schema_migrations ORDER BY version",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}

	applied := make(map[string]time.Time)
	for _, row := range rows {
		version := row["version"].(string)
		appliedAt := row["applied_at"].(time.Time)
		applied[version] = appliedAt
	}

	return applied, nil
}

// applyMigration applies a single migration
func (s *Service) applyMigration(m Migration) error {
	ctx := context.Background()

	// Execute the up migration
	if err := s.db.Exec(ctx, m.UpSQL); err != nil {
		return fmt.Errorf("failed to execute up migration: %w", err)
	}

	// Record the migration in schema_migrations
	if err := s.db.Exec(
		ctx,
		"INSERT INTO schema_migrations (version) VALUES ($1)",
		m.Version,
	); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return nil
}

// rollbackMigration rolls back a single migration
func (s *Service) rollbackMigration(m Migration) error {
	ctx := context.Background()

	// Execute the down migration
	if err := s.db.Exec(ctx, m.DownSQL); err != nil {
		return fmt.Errorf("failed to execute down migration: %w", err)
	}

	// Remove the migration from schema_migrations
	if err := s.db.Exec(
		ctx,
		"DELETE FROM schema_migrations WHERE version = $1",
		m.Version,
	); err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	return nil
}
