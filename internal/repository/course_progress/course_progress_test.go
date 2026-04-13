package courseprogress

import (
	"os"
	"fmt"
	"context"
	"testing"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"	

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	testDSN string
	testDB *gorm.DB
	pgContainer testcontainers.Container
)

const (
	BASE_POSTGRES_PORT = "5432"

	TEST_POSTGRES_PORT = "5444"
	TEST_POSTGRES_USER = "testProgress"
	TEST_POSTGRES_PASSWORD = "progress"
	TEST_POSTGRES_DB = "progressDB"
	TEST_POSTGRES_HOST = "localhost"

	MIGRATIONS_FOLDER = "migrations"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	pgContainer, err = testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image: "postgres:latest",
				ExposedPorts: []string{TEST_POSTGRES_PORT+"/tcp"},
				Env: map[string]string{
					"POSTGRES_USER": TEST_POSTGRES_USER,
					"POSTGRES_PASSWORD": TEST_POSTGRES_PASSWORD,
					"POSTGRES_DB": TEST_POSTGRES_DB,
				},
				WaitingFor: wait.ForListeningPort(BASE_POSTGRES_PORT+"/tcp"),
			},
			Started: true,
		},
	)
	if err != nil {
		panic(fmt.Errorf("start postgres container: %w", err))
	}

	testDB, err = gorm.Open(postgres.Open(getTestDSN()), &gorm.Config{})		
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		panic(fmt.Errorf("open db: %w", err))	
	}

	if err := runMigrations(getTestDSN()); err != nil {
		_ = pgContainer.Terminate(ctx)
		panic(fmt.Errorf("run migrations: %w", err))	
	}

	code := m.Run()

	_ = pgContainer.Terminate(ctx)

	os.Exit(code)
}

func getTestDSN() string {
	return "postgres://" + TEST_POSTGRES_USER + 
	":" + TEST_POSTGRES_PASSWORD + "@" + TEST_POSTGRES_HOST + ":" + 
	TEST_POSTGRES_PORT + "/" + TEST_POSTGRES_DB + "?sslmode=disable"
}

func runMigrations(dsn string) error {
	migrationsPath, err := filepath.Abs(MIGRATIONS_FOLDER)
	if err != nil {
		return err
	}

	m, err := migrate.New("file://" + migrationsPath, dsn)
	if err != nil {
		return err
	}

	defer func() {
		_, _ = m.Close()
	}()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func TestCreateCheckPoint()
func TestReadCheckPointsByUserID()

