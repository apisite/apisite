package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"html/template"
	"net/http"
)

// ErrRedirect is an error returned when page needs to be redirected
var ErrRedirect = errors.New("Abort with redirect")

// Meta holds page attributes
type Meta struct {
	Title        string
	JS           []string
	CSS          []string
	contentType  string
	status       int
	error        error
	errorMessage string
	layout       string
	location     string
}

// SetError sets error by template engine
// Not for use in templates (see Raise)
func (p *Meta) SetError(e error) {
	p.error = e
	if p.status == http.StatusOK {
		// error status was not set
		p.status = http.StatusInternalServerError
	}
}

func (p Meta) Error() error        { return p.error }
func (p Meta) Layout() string      { return p.layout }
func (p Meta) ContentType() string { return p.contentType }
func (p Meta) Status() int         { return p.status }
func (p Meta) Location() string    { return p.location }

// ErrorMessage returns internal or template error
func (p Meta) ErrorMessage() string {
	if p.errorMessage != "" {
		return p.errorMessage
	}
	if p.error == nil {
		return ""
	}
	return p.error.Error()
}

// SetLayout - set page layout
func (p *Meta) SetLayout(name string) (string, error) {
	p.layout = name
	return "", nil
}

// SetTitle - set page title
func (p *Meta) SetTitle(name string) (string, error) {
	p.Title = name
	return "", nil
}

// AddJS - add .js file to scripts list
func (p *Meta) AddJS(file string) (string, error) {
	p.JS = append(p.JS, file)
	return "", nil
}

// AddCSS - add .css file to styles list
func (p *Meta) AddCSS(file string) (string, error) {
	p.JS = append(p.CSS, file)
	return "", nil
}

// SetContentType - set page content type
func (p *Meta) SetContentType(name string) (string, error) {
	p.contentType = name
	return "", nil
}

// Raise - abort template processing (if given) and raise error
func (p *Meta) Raise(status int, abort bool, message string) (string, error) {
	p.status = status
	p.errorMessage = message // TODO: pass it via error only
	if abort {
		return "", errors.New(message)
	}
	return "", nil
}

// RedirectFound - abort template processing and return redirect with StatusFound status
func (p *Meta) RedirectFound(uri string) (string, error) {
	p.status = http.StatusFound
	p.location = uri
	return "", ErrRedirect // TODO: Is there a way to pass status & title via error?
}

// SetProtoFuncs appends function templates and not related to request functions to funcs
func SetProtoFuncs(funcs template.FuncMap) {
	funcs["request"] = func() interface{} { return nil }
	funcs["param"] = func(key string) string { return "" }
}

// SetRequestFuncs appends funcs which return real data inside request processing
func SetRequestFuncs(funcs template.FuncMap, ctx *gin.Context) {
	funcs["request"] = func() interface{} { return ctx.Request }
	funcs["param"] = func(key string) string { return ctx.Param(key) }
}
