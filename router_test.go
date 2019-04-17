package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	cfg := &Config{}
	log := setupLog()
	r, err := setupRouter(cfg, log)
	assert.NotNil(t, err)
	//	require.NoError(t, err)
	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	if err == nil {
		r.ServeHTTP(resp, req)

	}
	assert.Equal(t, resp.Body.String(), "")
}
