package main

import (
	//	"fmt"
	"html/template"
	"net/url"
	"strconv"

	"github.com/acoshift/paginate"
	"github.com/spf13/cast"
)

// Page holds pagination functionality
type Page struct {
	ID      int64 // Page number
	IsPrev  bool  // Is it a 'prev' link
	IsNext  bool  // Is it a 'last' link
	Allowed bool
	Href    string // Page href

	//	IsCurrent   bool // Is this a current page
	//	IsFirst     bool // Is it a first page
	//	IsLast      bool // Is it a last page
}

func valWithDefaultInt64(args url.Values, key string, valDefault int64) int64 {
	a, ok := args[key]
	if !ok {
		return valDefault
	}
	if len(a) == 0 {
		return valDefault
	}
	return cast.ToInt64(a[0])
}

// pager returns array with pagination links
func pager(args url.Values, count interface{}, argPrefix string, rowsMax, around, edge int64) *[]Page {
	itemCount := cast.ToInt64(count)

	rowLimit := valWithDefaultInt64(args, argPrefix+"lim", rowsMax)
	rowOffset := valWithDefaultInt64(args, argPrefix+"off", 0)

	pn := paginate.FromLimitOffset(rowLimit, rowOffset, itemCount)
	pageCurrent := rowOffset/rowLimit + 1
	p := []Page{}
	p = append(p, Page{IsPrev: true, Allowed: pn.CanPrev(), ID: pn.Prev(), Href: ""})
	for _, pi := range pn.Pages(around, edge) {
		p = append(p, Page{Allowed: pi != pageCurrent, ID: pi, Href: ""})
	}
	p = append(p, Page{IsNext: true, Allowed: pn.CanNext(), ID: pn.Next(), Href: ""})

	if rowLimit != rowsMax {
		args.Set(argPrefix+"lim", cast.ToString(rowLimit))
	} else {
		args.Del(argPrefix + "lim")
	}

	for i := range p {
		if p[i].ID > 1 {
			args.Set(argPrefix+"off", cast.ToString(rowLimit*(p[i].ID-1)))
		} else {
			args.Del(argPrefix + "off")
		}
		p[i].Href = "?" + args.Encode()
	}
	return &p
}

func interator(min, max int) []int {
	var ret []int
	if max == 0 {
		return ret
	}
	for i := min; i <= max; i++ {
		ret = append(ret, i)
	}
	return ret
}

func add(a, b int) int {
	return a + b
}

func atoi(a string, def int) int {
	if a == "" {
		return def
	}
	rv, err := strconv.Atoi(a)
	if err != nil {
		return def
	}
	return rv

}

// SetSimpleFuncs registers all prevoious funcs in given FuncMap
// This is not used in templates
func SetSimpleFuncs(funcs template.FuncMap) {
	funcs["add"] = add
	funcs["atoi"] = atoi
	funcs["pager"] = pager
	funcs["interator"] = interator
	funcs["bool"] = func(a bool) string {
		if a {
			return "+"
		}
		return "-"
	}
	funcs["HTML"] = func(s string) template.HTML {
		return template.HTML(s)
	}
	funcs["version"] = func() string {
		return version
	}
	funcs["ref"] = func(a *interface{}) interface{} {
		return *a
	}
}
