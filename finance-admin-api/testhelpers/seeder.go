package testhelpers

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	dbname   = "test_db"
	user     = "test_user"
	password = "test_password"
)

// Seeder contains a test database and API client for seeding data in tests
type Seeder struct {
	DbInstance *TestDatabase
	FinanceHub *TestFinanceHub
}

// NewSeeder creates a new Seeder instance
func NewSeeder() *Seeder {
	ctx := context.Background()

	os.Setenv("TESTCONTAINERS_LOG_LEVEL", "DEBUG")
	testcontainers.Logger = log.New(os.Stdout, "testcontainers: ", log.LstdFlags)

	db := InitDb(ctx)
	log.Printf("Running at: %s", db.Address)

	migrator, err := migrateDb(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Migrator started")

	// wait for the migrator to finish
	for {
		if migrator.IsRunning() {
			break
		}
		time.Sleep(1 * time.Second)
	}

	log.Printf("Migrator finished")

	err = db.DB.Snapshot(ctx, postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Snapshot taken")

	// Start the API startDbContainer
	fh, err := initFinanceHub(ctx, db.Address)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Finance hub started")

	return &Seeder{
		DbInstance: db,
		FinanceHub: fh,
	}
}

// Restore restores the DB to the snapshot backup and re-establishes the connection
func (s *Seeder) Restore() {
	err := s.DbInstance.DB.Restore(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

// TestDatabase is a test utility containing a fully-migrated Postgres instance. To use this, run InitDb within a TestMain
// function and use the DbInstance to interact with the database as needed (e.g. to insert data prior to testing).
// Ensure to run TearDown at the end of the tests to clean up.
type TestDatabase struct {
	Address string
	DB      *postgres.PostgresContainer
}

func InitDb(ctx context.Context) *TestDatabase {
	dbContainer, err := startDbContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	connString, err := dbContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return &TestDatabase{
		DB:      dbContainer,
		Address: connString,
	}
}

func startDbContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	_, b, _, _ := runtime.Caller(0)
	testPath := filepath.Dir(b)

	return postgres.Run(
		ctx,
		"docker.io/postgres:13-alpine",
		postgres.WithDatabase(dbname),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		postgres.WithInitScripts(testPath+"/public_schema.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
}

func migrateDb(ctx context.Context, db *TestDatabase) (testcontainers.Container, error) {
	connString, _ := db.DB.ConnectionString(ctx)
	log.Println(connString)

	ip, _ := db.DB.ContainerIP(ctx)

	port, err := db.DB.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-finance-migration:latest",
			Env: map[string]string{
				"DB_USER":       user,
				"DB_PASSWORD":   password,
				"DB_CONNECTION": "http://" + ip + ":" + port.Port(),
				"DB_NAME":       dbname,
				"DB_SCHEMA":     "supervision_finance",
			},
			Cmd: []string{"up"},
			WaitingFor: wait.ForSQL("5432/tcp", "postgres", func(h string, p nat.Port) string {
				return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?search_path=supervision_finance", user, password, h, p, dbname)
			}).WithStartupTimeout(5 * time.Minute),
		},
		Started: true,
	})
}

func (s *Seeder) TearDown() {
	_ = s.DbInstance.DB.Terminate(context.Background())
	_ = (*s.FinanceHub.container).Terminate(context.Background())
}

func (s *Seeder) GetConn() TestConn {
	conn, err := pgxpool.New(context.Background(), s.DbInstance.Address)
	if err != nil {
		log.Fatal(err)
	}
	return TestConn{conn}
}

type TestConn struct {
	Conn *pgxpool.Pool
}

func (c TestConn) Exec(ctx context.Context, s string, i ...interface{}) (pgconn.CommandTag, error) {
	return c.Conn.Exec(ctx, s, i...)
}

func (c TestConn) Query(ctx context.Context, s string, i ...interface{}) (pgx.Rows, error) {
	return c.Conn.Query(ctx, s, i...)
}

func (c TestConn) QueryRow(ctx context.Context, s string, i ...interface{}) pgx.Row {
	return c.Conn.QueryRow(ctx, s, i...)
}

func (c TestConn) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.Conn.BeginTx(ctx, pgx.TxOptions{})
}

func (c TestConn) SeedData(data ...string) {
	ctx := context.Background()
	for _, d := range data {
		_, err := c.Exec(ctx, d)
		if err != nil {
			log.Fatal("Unable to seed data with db connection: " + err.Error())
		}
	}
}

type TestFinanceHub struct {
	BaseURL    string
	container  *testcontainers.Container
	HTTPClient *http.Client
}

func initFinanceHub(ctx context.Context, connString string) (*TestFinanceHub, error) {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-finance-hub:latest",
			ExposedPorts: []string{"8080/tcp"},
			Env: map[string]string{
				"POSTGRES_CONN":     connString,
				"POSTGRES_USER":     user,
				"POSTGRES_PASSWORD": password,
				"POSTGRES_DB":       dbname,
				"DB_CONN_STRING":    connString,
			},
			WaitingFor: wait.ForHTTP("/finance/health-check").WithPort("8080/tcp"),
		},
		Started: true,
	})
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := container.MappedPort(ctx, "8080/tcp")
	if err != nil {
		return nil, err
	}

	baseURL := fmt.Sprintf("https://%s:%s", host, port.Port())
	return &TestFinanceHub{
		BaseURL:    baseURL,
		container:  &container,
		HTTPClient: &http.Client{},
	}, nil
}
