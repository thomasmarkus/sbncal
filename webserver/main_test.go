package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	router := setupRouter()

	//Test the index
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	//Test the wedstrijden api call
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/wedstrijden", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	//Test non-existing page
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/non-existing-page", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}
