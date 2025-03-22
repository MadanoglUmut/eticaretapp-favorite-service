package repositories

import (
	"context"
	"favorite_service/pkg/psql"
	"fmt"
	"os"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

var db *gorm.DB

var ctx context.Context = context.Background()

func TestMain(m *testing.M) {

	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"
	dbHost := "0.0.0.0"

	defaultPort := nat.Port("5432/tcp")

	postgresContainer, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{defaultPort.Port()},
			Env: map[string]string{
				"POSTGRES_USER":     dbUser,
				"POSTGRES_PASSWORD": dbPassword,
				"POSTGRES_DB":       dbName,
			},
			WaitingFor: wait.ForAll(
				wait.ForLog("database system is ready to accept connections"),
				wait.ForListeningPort("5432/tcp"),
			),
		},
		Started: true,
	})

	if err != nil {
		return
	}

	port, err := postgresContainer.MappedPort(ctx, defaultPort)

	if err != nil {
		return
	}

	fmt.Println("Default Port:", port)

	fileCreate, err := os.ReadFile("../../psql/create_tables.sql")

	if err != nil {
		return
	}

	fileFill, err := os.ReadFile("../../psql/fill_tables.sql")

	db = psql.Connect(dbHost, dbUser, dbPassword, dbName, port.Port())

	err = db.Exec(string(fileCreate)).Error

	if err != nil {
		return
	}

	err = db.Exec(string(fileFill)).Error

	if err != nil {
		return
	}
	os.Exit(m.Run())
}
