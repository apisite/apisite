package tplfunc

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestNewMeta(t *testing.T) {

	m := NewMeta(http.StatusNotImplemented, "text/plain")
	assert.Equal(t, http.StatusNotImplemented, m.Status())
	assert.Equal(t, "", m.ErrorMessage(), "no error")

	msg := "Interrupted"
	_, err := m.Raise(http.StatusNotAcceptable, true, msg)
	assert.Equal(t, msg, err.Error(), "err message")
	assert.Equal(t, msg, m.ErrorMessage(), "error message")
	assert.Equal(t, http.StatusNotAcceptable, m.Status(), "status")
}

func TestRaiseNoAbort(t *testing.T) {

	m := Meta{}

	msg := "Not Interrupted"
	_, err := m.Raise(http.StatusNotAcceptable, false, msg)
	require.NoError(t, err, "store error without ebort")
	assert.Equal(t, msg, m.ErrorMessage(), "error message")
	assert.Equal(t, http.StatusNotAcceptable, m.Status(), "status")
}

func TestErrorMessageExternal(t *testing.T) {
	m := Meta{}
	msg := "ExternalError"
	e := errors.New(msg)
	m.SetError(e)
	assert.Equal(t, msg, m.ErrorMessage())
}

func TestRedirectFound(t *testing.T) {
	m := Meta{}
	_, err := m.RedirectFound("/redirect")

	assert.Equal(t, "Abort with redirect", err.Error())
	assert.Equal(t, "/redirect", m.Location())
}
