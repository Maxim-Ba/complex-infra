package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"go-auth/internal/app"
	"go-auth/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	db *sql.DB
}

func New(cfg app.AppConfig) *DB {
	var c = cfg.GetConfig()

	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB,
	)
	pgxCfg, err := pgx.ParseConfig(connString)
	if err != nil {
		panic(err.Error())
	}

	db := stdlib.OpenDB(*pgxCfg)

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		panic(err.Error())
	}
	checkDerectory(c)
	if err := applyMigrations(db, c.MigrationPath); err != nil {
		panic(err.Error())
	}
	
	fmt.Println("Postgres connection created", " host: ", c.PostgresHost, " port: ", c.PostgresPort)
	return &DB{db: db}
}


func (db *DB) Close() {
	db.db.Close()
}
func (db *DB) GetConnection() *sql.DB {
	return db.db
}
func applyMigrations(db *sql.DB, migrationPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationPath),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not apply migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("could not get migration version: %w", err)
	}

	fmt.Printf("Migrations applied successfully. Current version: %d, dirty: %v", version, dirty)
	slog.Info("Migrations applied successfully")
	return nil
}
func checkDerectory(cfg *config.Config){
	dir, err := os.Getwd()
    if err != nil {
        panic(fmt.Sprintf("failed to get current directory: %v", err))
    }
    fmt.Printf("Current working directory: %s\n", dir)

    // Проверка существования папки миграций
    migrationPath := cfg.GetConfig().MigrationPath
    absPath, _ := filepath.Abs(migrationPath)
    fmt.Printf("Migration path (abs): %s\n", absPath)

    if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
        panic(fmt.Sprintf("migrations directory does not exist: %s", migrationPath))
    }
}
