package main

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/matryer/is"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestUserRepository(t *testing.T) {
	// skip in short mode
	if testing.Short() {
		return
	}

	is := is.New(t)

	// Setup database
	dbContainer, connPool, err := SetupTestDatabase()
	if err != nil {
		t.Error(err)
	}
	defer dbContainer.Terminate(context.Background())

	// Create user repository
	repository := NewUserRepository(connPool)

	// Run tests against db
	t.Run("FindExistingUserByUsername", func(t *testing.T) {
		adminUser, err := repository.FindByUsername(
			context.Background(),
			"admin",
		)

		is.NoErr(err)
		is.Equal(adminUser.Username, "admin")
	})

}

func SetupTestDatabase() (testcontainers.Container, *pgxpool.Pool, error) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}
	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		return nil, nil, err
	}
	port, err := dbContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		return nil, nil, err
	}
	host, err := dbContainer.Host(context.Background())
	if err != nil {
		return nil, nil, err
	}

	dbURI := fmt.Sprintf("postgres://postgres:postgres@%v:%v/testdb", host, port.Port())
	err = MigrateDb(dbURI)
	if err != nil {
		return nil, nil, err
	}

	connPool, err := pgxpool.Connect(context.Background(), dbURI)
	if err != nil {
		return nil, nil, err
	}

	return dbContainer, connPool, err
}
