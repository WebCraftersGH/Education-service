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
	"github.com/google/uuid"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/WebCraftersGH/Education-service/pkg/logging"
	"github.com/WebCraftersGH/Education-service/internal/config"
	"github.com/WebCraftersGH/Education-service/internal/domain"
)

var (
	testDSN string
	testDB *gorm.DB
	pgContainer testcontainers.Container
	cfg config.Config
	logger logging.Logger
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

var (
	TEST_USER_UUID = uuid.New()
	TEST_USER_UUID2 = uuid.New()
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

	cfg, err = config.Load(".env")
	if err != nil {
		panic(fmt.Errorf("load config: %w", err))
	}

	logger, closer, err := logging.New(cfg.LogLevel)
	if err != nil {
		panic(fmt.Errorf("get logger: %w", err))
	}
	defer closer.Close()

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


//TestCreateCheckPoint checks the creation of checkpoints
func TestCreateCheckPoint(t *testing.T) {	

	tests := []struct{
		name string
		checkPoint domain.CheckPoint
		creationStatus bool
	}{
		{"first user with unique slug", domain.CheckPoint{
			ID: uuid.New(),
			UserID: TEST_USER_UUID,
			Slug: "GitCourse__Module-first",
		}, true},
		{"second user with unique slug", domain.CheckPoint{
			ID: uuid.New(),
			UserID: TEST_USER_UUID2,
			Slug: "GitCourse__Module-second",
		}, true},
		{"first user with second user slug", domain.CheckPoint{
			ID: uuid.New(),
			UserID: TEST_USER_UUID,
			Slug: "GitCourse__Module-second",
		}, true},
		// the user must have a unique slug
		{"first user with first slug", domain.CheckPoint{
			ID: uuid.New(),
			UserID: TEST_USER_UUID,
			Slug: "GitCourse__Module-first",
		}, false},
		{"first user with new slug", domain.CheckPoint{
			ID: uuid.New(),
			UserID: TEST_USER_UUID,
			Slug: "GitCourse__Module-third",
		}, true},
	}

	ctx := context.Background()
	repo := NewProgressRepo(testDB, logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCheckpoint, err := repo.CreateCheckPoint(ctx, tt.checkPoint)
			switch {
			case err != nil && tt.creationStatus == true:
				t.Errorf("failed to create checkpoint: %w", err)
			case err == nil && tt.creationStatus == true && gotCheckPoint.Slug != tt.checkPoint.Slug:
				t.Errorf("got checkpoint != result checkpoint; %s != %s", gotCheckPoint.Slug, tt.checkPoint.Slug)
			}
		})
	}
}

func TestReadCheckPointsByUserID(t *testing.T) {

	tests := []struct{
		name string
		limit int
		offset int
		userID uuid.UUID
		recordsCount int
	}{
		{"get first user checkpoints", 10, 0, TEST_USER_UUID, 3},
		{"get second user checkpoints", 10, 0, TEST_USER_UUID2, 1},
		{"get first user checkpoints with offset", 10, 1, TEST_USER_UUID, 2},
		{"get first user checkpoints with negative limit", -1, 0, TEST_USER_UUID, 3},
		{"get first user checkpoints with negative offset", 10, -1, TEST_USER_UUID, 3}	
	}

	ctx := context.Background()
	repo := NewProgressRepo(testDB, logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCheckPoints, err := repo.ReadCheckPointsByUserID(ctx, tt.userID, tt.limit, tt.offset)
			switch {
			case err != nil:
				t.Errorf("failed to get checkpoints list by user (%s): %w", TEST_USER_UUID.String(), err)
			case err == nil && len(gotCheckPoints) != tt.recordsCount:
				t.Errorf("len of got checkpoints != expected len; %d != %d", len(gotCheckPoints), tt.recordsCount)
			}
		})
	}
}

