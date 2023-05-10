package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	router := setupRouter()

	//Test the index
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, nil, err)
	assert.Equal(t, 200, w.Code)

	//Test importing wedstrijden
	wedstrijden, err = importWedstrijden("../wedstrijden.json")

	assert.Equal(t, nil, err)
	assert.Less(t, 0, len(wedstrijden))

	//Test the wedstrijden api call
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/api/wedstrijden", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, nil, err)
	assert.Equal(t, 200, w.Code)

	//Check if the response contains a valid wedstrijden array
	var wedstrijdenTest []Wedstrijd
	err = json.Unmarshal(w.Body.Bytes(), &wedstrijdenTest)

	assert.Equal(t, err, nil)
	assert.Equal(t, wedstrijden, wedstrijdenTest)

	//Test non-existing page
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/api/non-existing-page", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, nil, err)
	assert.Equal(t, 404, w.Code)
}
