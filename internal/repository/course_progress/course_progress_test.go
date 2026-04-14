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
	"github.com/WebCraftersGH/Education-service/internal/domain"
)

var (
	testDB *gorm.DB
	pgContainer testcontainers.Container
	logger logging.Logger
)

const (
	TEST_POSTGRES_PORT = "5432/tcp"
	TEST_POSTGRES_USER = "testProgress"
	TEST_POSTGRES_PASSWORD = "progress"
	TEST_POSTGRES_DB = "progressDB"
	TEST_POSTGRES_HOST = "localhost"

	MIGRATIONS_FOLDER = "../../../migrations"
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
				ExposedPorts: []string{TEST_POSTGRES_PORT},
				Env: map[string]string{
					"POSTGRES_USER": TEST_POSTGRES_USER,
					"POSTGRES_PASSWORD": TEST_POSTGRES_PASSWORD,
					"POSTGRES_DB": TEST_POSTGRES_DB,
				},
				WaitingFor: wait.ForListeningPort(TEST_POSTGRES_PORT),
			},
			Started: true,
		},
	)
	if err != nil {
		panic(fmt.Errorf("start postgres container: %v", err))
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		panic(err)
	}

	port, err := pgContainer.MappedPort(ctx, TEST_POSTGRES_PORT)
	if err != nil {
		panic(err)
	}

	testDB, err = gorm.Open(postgres.Open(getTestDSN(host, port.Port())), &gorm.Config{})		
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		panic(fmt.Errorf("open db: %v", err))	
	}

	if err := runMigrations(getTestDSN(host, port.Port())); err != nil {
		_ = pgContainer.Terminate(ctx)
		panic(fmt.Errorf("run migrations: %v", err))	
	}

	logger, _, err = logging.New("info")
	if err != nil {
		panic(fmt.Errorf("get logger: %v", err))
	}

	code := m.Run()

	_ = pgContainer.Terminate(ctx)

	os.Exit(code)
}

func getTestDSN(host, port string) string {
	return "postgres://" + TEST_POSTGRES_USER + 
	":" + TEST_POSTGRES_PASSWORD + "@" + host + ":" + 
	port + "/" + TEST_POSTGRES_DB + "?sslmode=disable"
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

//clearDB helper func for clear database
func clearDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	err := db.Exec(`
	TRUNCATE TABLE checkpoints RESTART IDENTITY CASCADE;
	`).Error
	
	if err != nil {
		t.Fatalf("clear db: %v", err)
	}
}

func seedCheckpoints(
	t *testing.T, 
	repo *progressRepo, 
	ctx context.Context, 
	checkpoints []domain.CheckPoint,
) {
	t.Helper()

	for _, cp := range checkpoints {
		if _, err := repo.CreateCheckPoint(ctx, cp); err != nil {
			t.Fatalf("seed checkpoint: %v", err)
		}
	}
}

//TODO Данный тест боится параллельности. Зависит от порядка внутренних тестов.
//TestCreateCheckPoint checks the creation of checkpoints
func TestCreateCheckPoint(t *testing.T) {	

	clearDB(t, testDB)

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
				t.Errorf("failed to create checkpoint: %v", err)
			case err == nil && tt.creationStatus == true && gotCheckpoint.Slug != tt.checkPoint.Slug:
				t.Errorf("got slug %s, want %s", gotCheckpoint.Slug, tt.checkPoint.Slug)
			case err == nil && !tt.creationStatus:
				t.Errorf("expected error, got nil")
			}
		})
	}
}

func TestReadCheckPointsByUserID(t *testing.T) {

	clearDB(t, testDB)

	baseData := []domain.CheckPoint{
		{
			ID:     uuid.New(),
			UserID: TEST_USER_UUID,
			Slug:   "GitCourse__Module-first",
		},
		{
			ID:     uuid.New(),
			UserID: TEST_USER_UUID,
			Slug:   "GitCourse__Module-second",
		},
		{
			ID:     uuid.New(),
			UserID: TEST_USER_UUID,
			Slug:   "GitCourse__Module-third",
		},
		{
			ID:     uuid.New(),
			UserID: TEST_USER_UUID2,
			Slug:   "GitCourse__Module-second",
		},
	}

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
		{"get first user checkpoints with negative offset", 10, -1, TEST_USER_UUID, 3},
	}

	ctx := context.Background()
	repo := NewProgressRepo(testDB, logger)

	seedCheckpoints(t, repo, ctx, baseData)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCheckPoints, err := repo.ReadCheckPointsByUserID(ctx, tt.userID, tt.limit, tt.offset)
			switch {
			case err != nil:
				t.Errorf("failed to get checkpoints list by user (%s): %v", tt.userID.String(), err)
			case err == nil && len(gotCheckPoints) != tt.recordsCount:
				t.Errorf("len of got checkpoints != expected len; %d != %d", len(gotCheckPoints), tt.recordsCount)
			}
		})
	}
}

