package test

import (
	"context"
	"genuinebnt/newsletter/config"
	lib "genuinebnt/newsletter/internal"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
)

func spawnDB(t testing.TB) *pgxpool.Pool {
	t.Helper()
	config, err := config.GetConfiguration()
	if err != nil {
		t.Fatal("Failed to read configuration: ", err)
	}

	name, err := uuid.NewRandom()
	if err != nil {
		t.Fatal(err)
	}

	config.Database.Name = name.String()
	return configureDatabase(t, config)
}

func configureDatabase(t testing.TB, config *config.Config) *pgxpool.Pool {
	t.Helper()
	connectionString := config.Database.ConnnectionStringWithoutDB()

	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		t.Fatal("Failed to connect to database: ", err)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "CREATE DATABASE \""+config.Database.Name+"\"")
	if err != nil {
		t.Fatal("Failed to create database: ", err)
	}

	dbpool, err := pgxpool.New(context.Background(), config.Database.ConnectionString())
	if err != nil {
		t.Fatal("Failed to connect to database: ", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	migrationsDir := filepath.Join(dir, "../migrations")

	//var embedMigrations embed.FS
	goose.SetBaseFS(nil)

	db := stdlib.OpenDBFromPool(dbpool)
	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatal(err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		t.Fatal(err)
	}

	return dbpool
}

func TestHealthCheck(t *testing.T) {
	dbpool := spawnDB(t)
	defer dbpool.Close()

	server := httptest.NewServer(lib.Server(dbpool))
	defer server.Close()

	resp, err := http.Get(server.URL + "/health_check")
	if err != nil {
		t.Fatal("Failed to execute http request with err", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, resp.ContentLength, int64(0))
}

func TestSubscribe(t *testing.T) {
	t.Run("Subscriber returns 200 for a valid form data", func(t *testing.T) {
		dbpool := spawnDB(t)
		defer dbpool.Close()

		server := httptest.NewServer(lib.Server(dbpool))
		defer server.Close()

		body := "name=genuine%20basil%20nt&email=genuinebnt%40gmail.com"
		resp, err := http.Post(server.URL+"/subscriptions", "application/x-www-form-urlencoded", strings.NewReader(body))
		if err != nil {
			t.Fatal("Failed to execute http request with err", err)
		}

		var name string
		var email string
		err = dbpool.QueryRow(context.Background(), "SELECT email, name from subscriptions;").Scan(&email, &name)
		if err != nil {
			t.Fatal("Failed to fetch subscriptions: ", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		assert.Equal(t, name, "genuine basil nt")
		assert.Equal(t, email, "genuinebnt@gmail.com")
	})

	t.Run("Subscriber returns 400 when data is missing", func(t *testing.T) {
		dbpool := spawnDB(t)
		defer dbpool.Close()

		server := httptest.NewServer(lib.Server(dbpool))
		defer server.Close()

		var testCases = []struct {
			input string
			err   string
		}{
			{
				input: "name=genuine%20basil%20nt",
				err:   "missing required field: email",
			},
			{
				input: "email=genuinebnt%40gmail.com",
				err:   "missing required field: name",
			},
			{
				input: "",
				err:   "missing required field: name and email",
			},
		}

		for _, testCase := range testCases {
			resp, err := http.Post(server.URL+"/subscriptions", "application/x-www-form-urlencoded", strings.NewReader(testCase.input))
			if err != nil {
				t.Fatal("Failed to execute http request with err", err)
			}

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "The API did not fail with 400 status code when payload was "+testCase.err)
		}
	})
}
