package initialize

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"time"
	"urlShortener/internal/utils"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewClient(ctx context.Context, maxAttempts int, config *Config) (*DB, error) {
	var pool *pgxpool.Pool
	var err error
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		config.PGUser,
		config.PGPassword,
		config.PGHost,
		config.PGPort,
		config.PGDatabase)

	fmt.Println(config)

	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		return nil, err
	}

	dbInstance := &DB{
		Pool: pool,
	}
	//if err := dbInstance.RunMigrations(pool); err != nil {
	//	pool.Close()
	//	return nil, err
	//}

	return dbInstance, nil
}

func (d *DB) RunMigrations(logger *zap.Logger) error {
	db := stdlib.OpenDBFromPool(d.Pool)

	if err := goose.SetDialect("pgx"); err != nil {
		logger.Error("Error setting goose dialect", zap.Error(err))
		return err
	}

	migrationsDir := "./schema"

	if err := goose.Up(db, migrationsDir); err != nil {
		logger.Error("Error running migrations", zap.Error(err))
		return err
	}

	logger.Info("Migrations successfully applied")
	//fmt.Println("Migrations applied successfully!")
	return nil
}
