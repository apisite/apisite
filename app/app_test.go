package app

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	// Save original args
	a := os.Args

	tests := []struct {
		name string
		code int
		args []string
	}{
		{"Help", 3, []string{"-h"}},
		{"UnknownFlag", 2, []string{"--unknown"}},
		{"UnknownPort", 1, []string{"--http_addr", ":xx"}},
	}
	for _, tt := range tests {
		os.Args = append([]string{a[0]}, tt.args...)
		var c int
		Run(func(code int) { c = code })
		assert.Equal(t, tt.code, c, tt.name)
	}

	// Restore original args
	os.Args = a
}

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
