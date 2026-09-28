package db

import (
	"fmt"

	"github.com/Snehaaa-Kri/job-apply-automator/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewPostgresDB creates and returns a GORM DB connection.
// Node analogy: creating a pg.Pool or Sequelize instance — connect once, reuse everywhere.
func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	// Build DSN (connection string)
	dsn := cfg.Database.DSN()

	// Log level: show SQL queries in dev, silence in prod
	logLevel := gormlogger.Info
	if cfg.App.Env == "production" {
		logLevel = gormlogger.Error
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Connection pool settings
	// Node analogy: pool: { max: 10, min: 2 } in pg config
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)     // max simultaneous connections
	sqlDB.SetMaxIdleConns(10)     // idle connections kept alive

	return db, nil
}
