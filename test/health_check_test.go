package test

import (
	lib "genuinebnt/newsletter/internal"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	server := httptest.NewServer(lib.Server())
	defer server.Close()

	resp, err := http.Get(server.URL + "/health_check")
	if err != nil {
		log.Fatalln("Failed to execute http request with err", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int64(0), resp.ContentLength)
}
