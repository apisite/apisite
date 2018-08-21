package main

import (
	//	"fmt"
	"html/template"
	"net/url"
	"strconv"

	"github.com/spf13/cast"
)

type Page struct {
	ID          int    // Page number
	IsCurrent   bool   // Is this a current page
	IsFirst     bool   // Is it a first page
	IsLast      bool   // Is it a last page
	IsPrev      bool   // Is it 'prev' link
	IsNext      bool   // 'last' link
	SpaceBefore bool   // There is a space before this item
	SpaceAfter  bool   // There is a space before this item
	Href        string // Page href
}

func valWithDefaultInt(args *url.Values, key string, valDefault int) int {
	a, ok := (*args)[key]
	if !ok {
		return valDefault
	}
	if len(a) == 0 {
		return valDefault
	}
	return cast.ToInt(a[0])
}

//
func pager(args url.Values, count *interface{}, argPrefix string, rowsMax, blockLimit int) (*[]Page, error) {
	itemCount := cast.ToInt(*count)

	rowLimit := valWithDefaultInt(&args, argPrefix+"lim", rowsMax)
	//rowOffset := valWithDefaultInt(&args, argPrefix+"off", 0)

	pageCount := 0
	if rowLimit > 0 {
		pageCount = int(itemCount / rowLimit)
		if itemCount%rowLimit > 0 {
			pageCount++
		}
	}

	pages := []Page{}
	/*
		max := pages - 1

		prev := "#"
		if pageNum == 1 {
			prev = "../" + prefix + "/"
		} else if pageNum > 1 {
			prev = strconv.Itoa(pageNum - 1)
		}
		next := "#"
		if pageNum < max {
			next = strconv.Itoa(pageNum + 1)
		}

		// ------------------------
		p := Pager{
			Count:   pages,
			Max:     max,
			Current: pageNum,
			Enabled: (pages > 0),
			IsFirst: (pageNum == 0),
			IsLast:  (pageNum == max),
			Prev:    prev,
			Next:    next,
		}
	*/
	return &pages, nil
}

/*
func pages(a string, def int) int {

	q := req.URL.Query()
	q.Add("api_key", "key_from_environment_or_flag")
	q.Set("api_key", "foo & bar")
	q.Del("api_key")
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())
}
*/
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

func appendFuncs(funcs template.FuncMap) {
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
}
