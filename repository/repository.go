package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type PostgresCfg struct {
	Host        string
	Port        string
	PG_USER     string
	PG_PASSWORD string
	PG_DB       string
	SSL_MODE    string
}

type postgres struct {
	db      *pgxpool.Pool
	querier Querier
}

var (
	pgInstance *postgres
	pgOnce     sync.Once
	initErr    error
)

type Repository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (int32, error)
	DeleteUserByID(ctx context.Context, id int32) (int64, error)
	GetUserByID(ctx context.Context, id int32) (User, error)
	UpdateUserPwd(ctx context.Context, arg UpdateUserPwdParams) (User, error)
	IsUserValid(ctx context.Context, arg string) (User, error)
	Close()
}

func NewPostgresDB(ctx context.Context, cfg PostgresCfg) (Repository, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.PG_USER,
		cfg.PG_PASSWORD,
		cfg.PG_DB,
		cfg.SSL_MODE,
	)

	pgOnce.Do(func() {

		dbpool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			logrus.Errorf("Unable to create pgxpool, %v", err)
			initErr = err
			return
		}

		if err := dbpool.Ping(ctx); err != nil {
			logrus.Errorf("Unable to pind db pool, %v", err)
			initErr = err
			return
		}

		queries := New(dbpool)

		pgInstance = &postgres{dbpool, queries}

		logrus.Info("Successfully initialized Postgres connection pool")

		migrateDSN := fmt.Sprintf("pgx5://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.PG_USER, cfg.PG_PASSWORD, cfg.Host, cfg.Port, cfg.PG_DB, cfg.SSL_MODE)

		m, err := migrate.New("file://migrations", migrateDSN+"&x-migrations-table=order_migrations")
		if err != nil {

			logrus.Errorf("Unable to find migrations, %v", err)
			initErr = err
			return

		}

		if err := m.Up(); err != nil {
			if err != migrate.ErrNoChange {
				logrus.Errorf("Unable to apply migrations, %v", err)
				initErr = err
				return
			}
		}

		logrus.Info("Migrations successfully applied")
	})

	if initErr != nil {
		return nil, initErr
	}

	if pgInstance == nil {
		return nil, initErr
	}

	return pgInstance, nil
}

func (r *postgres) CreateUser(ctx context.Context, arg CreateUserParams) (int32, error) {
	return r.querier.CreateUser(ctx, arg)
}

func (r *postgres) DeleteUserByID(ctx context.Context, id int32) (int64, error) {
	return r.querier.DeleteUserByID(ctx, id)
}

func (r *postgres) GetUserByID(ctx context.Context, id int32) (User, error) {
	return r.querier.GetUserByID(ctx, id)
}

func (r *postgres) UpdateUserPwd(ctx context.Context, arg UpdateUserPwdParams) (User, error) {
	return r.querier.UpdateUserPwd(ctx, arg)
}

func (r *postgres) IsUserValid(ctx context.Context, arg string) (User, error) {
	return r.querier.IsUserValid(ctx, arg)
}

func (r *postgres) Close() {
	r.db.Close()
}
