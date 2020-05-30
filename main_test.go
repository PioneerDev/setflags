package main

import (
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"set-flags/models"
	"set-flags/routers"
	"strings"
	"testing"
)

var router = routers.InitRouter()


func TestPing(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestListFlags(t *testing.T) {
	models.InitDB()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/flags", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateFlag(t *testing.T) {
	models.InitDB()

	reader := strings.NewReader(`
		"payer_id":       "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		"task":           "test",
		"days":           10,
		"max_witness":    100,
		"asset_id":       "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		"amount":         0.4,
		"times_achieved": 365,
		"status":         "paid",
	`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/flag", reader)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
