package test

import (
	"context"
	"genuinebnt/newsletter/config"
	lib "genuinebnt/newsletter/internal"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	config, err := config.GetConfiguration()
	if err != nil {
		log.Fatalln("Failed to read configuration: ", err)
	}

	connectionString := config.Database.ConnectionString()

	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		log.Fatalln("Failed to connect to database: ", err)
	}
	defer conn.Close(context.Background())

	server := httptest.NewServer(lib.Server(conn))
	defer server.Close()

	resp, err := http.Get(server.URL + "/health_check")
	if err != nil {
		log.Fatalln("Failed to execute http request with err", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, resp.ContentLength, int64(0))
}

func TestSubscribe(t *testing.T) {
	t.Run("Subscriber returns 200 for a valid form data", func(t *testing.T) {
		config, err := config.GetConfiguration()
		if err != nil {
			log.Fatalln("Failed to read configuration: ", err)
		}

		connectionString := config.Database.ConnectionString()

		conn, err := pgx.Connect(context.Background(), connectionString)
		if err != nil {
			log.Fatalln("Failed to connect to database: ", err)
		}
		defer conn.Close(context.Background())

		server := httptest.NewServer(lib.Server(conn))
		defer server.Close()

		body := "name=genuine%20basil%20nt&email=genuinebnt%40gmail.com"
		resp, err := http.Post(server.URL+"/subscriptions", "application/x-www-form-urlencoded", strings.NewReader(body))
		if err != nil {
			log.Fatalln("Failed to execute http request with err", err)
		}

		var name string
		var email string
		err = conn.QueryRow(context.Background(), "SELECT email, name from subscriptions;").Scan(&email, &name)
		if err != nil {
			log.Fatalln("Failed to fetch subscriptions: ", err)
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		assert.Equal(t, name, "genuine basil nt")
		assert.Equal(t, email, "genuinebnt@gmail.com")
	})

	t.Run("Subscriber returns 400 when data is missing", func(t *testing.T) {
		config, err := config.GetConfiguration()
		if err != nil {
			log.Fatalln("Failed to read configuration: ", err)
		}

		connectionString := config.Database.ConnectionString()

		conn, err := pgx.Connect(context.Background(), connectionString)
		if err != nil {
			log.Fatalln("Failed to connect to database: ", err)
		}
		defer conn.Close(context.Background())

		server := httptest.NewServer(lib.Server(conn))
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
				log.Fatalln("Failed to execute http request with err", err)
			}

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "The API did not fail with 400 status code when payload was "+testCase.err)
		}
	})
}
