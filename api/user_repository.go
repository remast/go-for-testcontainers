package main

import (
	"context"
	"embed"
	"errors"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//go:embed migrations
var migrations embed.FS

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID       uuid.UUID
	Username string
	Password string
}

// UserRepository is a SQL database repository for users
type UserRepository struct {
	connPool *pgxpool.Pool
}

// NewUserRepository creates a new SQL database repository for users
func NewUserRepository(connPool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		connPool: connPool,
	}
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	row := r.connPool.QueryRow(
		ctx,
		`SELECT user_id, password
		 FROM users 
		 WHERE username = $1`, username,
	)

	var (
		id       string
		password string
	)

	err := row.Scan(&id, &password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	user := &User{
		ID:       uuid.MustParse(id),
		Username: username,
		Password: password,
	}
	return user, nil
}

func MigrateDb(dbURI string) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, strings.Replace(dbURI, "postgres://", "pgx://", 1))
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
