package database

import (
	"database/sql"
	"fmt"

	defaultLogger "github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger/defaultlogger"

	"emperror.dev/errors"
	"github.com/glebarez/sqlite"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGorm(cfg *Options) (*gorm.DB, error) {
	if cfg.DBName == "" {
		return nil, errors.New("missing DBName in gorm configuration")
	}

	if cfg.UseSQLLite {
		db, err := createSQLLiteDB(cfg.DNS())
		return db, err
	}

	// InMemory doesn't work correctly with transactions - seems when we `Begin` a transaction on gorm.DB (with SQLLite in-memory) our previous gormDB before transaction will remove and the new gormDB with tx will go on the memory
	if cfg.UseInMemory {
		db, err := createInMemoryDB()
		return db, err
	}

	err := createPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
	)

	gormDb, err := gorm.Open(
		gormPostgres.Open(dataSourceName),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	return gormDb, nil
}

func createInMemoryDB() (*gorm.DB, error) {
	// https://gorm.io/docs/connecting_to_the_database.html#SQLite
	// https://github.com/glebarez/sqlite
	// https://www.connectionstrings.com/sqlite/
	db, err := gorm.Open(
		sqlite.Open(":memory:"),
		&gorm.Config{})

	return db, err
}

func createSQLLiteDB(dbFilePath string) (*gorm.DB, error) {
	// https://gorm.io/docs/connecting_to_the_database.html#SQLite
	// https://github.com/glebarez/sqlite
	// https://www.connectionstrings.com/sqlite/
	gormSQLLiteDB, err := gorm.Open(
		sqlite.Open(dbFilePath),
		&gorm.Config{})

	return gormSQLLiteDB, err
}

func createPostgresDB(cfg *Options) error {
	var db *sql.DB

	// we should choose a default database in the connection, but because we don't have a database yet we specify postgres default database 'postgres'
	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		"postgres",
		cfg.Password,
	)
	postgresGormDB, err := gorm.Open(
		gormPostgres.Open(dataSourceName),
		&gorm.Config{},
	)
	if err != nil {
		return err
	}

	db, err = postgresGormDB.DB()
	if err != nil {
		return err
	}

	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT 1 FROM pg_catalog.pg_database WHERE datname='%s'",
			cfg.DBName,
		),
	)
	if err != nil {
		return err
	}

	var exists int
	if rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return err
		}
	}

	if exists == 1 {
		return nil
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			defaultLogger.GetLogger().Error("Error closing database connection", err)
		}
	}(db)

	return nil
}
