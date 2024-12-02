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
	"path/filepath"
	"runtime"
	"time"
)

const (
	dbname   = "test_db"
	user     = "test_user"
	password = "test_password"
)

// TestDatabase is a test utility containing a fully-migrated Postgres instance. To use this, run InitDb within a TestMain
// function and use the DbInstance to interact with the database as needed (e.g. to insert data prior to testing).
// Ensure to run TearDown at the end of the tests to clean up.
type TestDatabase struct {
	Address string
	DB      *postgres.PostgresContainer
	API     *testcontainers.Container
}

// Restore restores the DB to the snapshot backup and re-establishes the connection
func (db *TestDatabase) Restore() {
	err := db.DB.Restore(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func InitDb() *TestDatabase {
	ctx := context.Background()

	dbContainer, err := startDbContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	connString, err := dbContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatal(err)
	}

	migrator, err := migrateDb(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}

	// wait for the migrator to finish
	for {
		if !migrator.IsRunning() {
			break
		}
		time.Sleep(1 * time.Second)
	}

	err = dbContainer.Snapshot(ctx, postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		log.Fatal(err)
	}

	// Start the API startDbContainer
	api, err := startApiContainer(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}

	return &TestDatabase{
		DB:      dbContainer,
		Address: connString,
		API:     &api,
	}
}

func startDbContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	_, b, _, _ := runtime.Caller(0)
	testPath := filepath.Dir(b)
	//basePath := filepath.Join(testPath, "../../..")

	return postgres.Run(
		ctx,
		"docker.io/postgres:13-alpine",
		postgres.WithDatabase(dbname),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		postgres.WithInitScripts(testPath+"/migrations/public_schema.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
}

func migrateDb(ctx context.Context, connString string) (testcontainers.Container, error) {
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-finance-migration:latest",
			Env: map[string]string{
				"DB_USER":       user,
				"DB_PASSWORD":   password,
				"DB_CONNECTION": connString,
				"DB_NAME":       dbname,
				"DB_SCHEMA":     "supervision_finance",
			},
			WaitingFor: wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
				return fmt.Sprintf("postgresql://%s:%s@%s/%s?search_path=supervision_finance", user, password, connString, dbname)
			}).WithStartupTimeout(5 * time.Second),
		},
		Started: true,
	})
}

func startApiContainer(ctx context.Context, connString string) (testcontainers.Container, error) {
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
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
}

func (db *TestDatabase) TearDown() {
	_ = db.DB.Terminate(context.Background())
}

func (db *TestDatabase) GetConn() TestConn {
	conn, err := pgxpool.New(context.Background(), db.Address)
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
