// Package tplfunc implements functions avalilable inside templates
package tplfunc

import (
	"github.com/gin-gonic/gin"
	"html/template"

	base "github.com/apisite/apitpl/ginapitpl/samplemeta"
)

// Meta holds template metadata
// This func is not for use in templates
type Meta struct {
	base.Meta
	JS  []string
	CSS []string
}

// NewMeta returns new initialised Meta struct
// This func is not for use in templates
func NewMeta(status int, ctype string) *Meta {
	m := base.NewMeta(status, ctype)
	return &Meta{Meta: *m}
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

// SetProtoFuncs appends function templates and not related to request functions to funcs
// This func is not for use in templates
func SetProtoFuncs(funcs template.FuncMap) {
	funcs["request"] = func() interface{} { return nil }
	funcs["param"] = func(key string) string { return "" }
}

// SetRequestFuncs appends funcs which return real data inside request processing
// This func is not for use in templates
func SetRequestFuncs(funcs template.FuncMap, ctx *gin.Context) {
	funcs["request"] = func() interface{} { return ctx.Request }
	funcs["param"] = func(key string) string { return ctx.Param(key) }
}
