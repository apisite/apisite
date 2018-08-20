package main

import (
	"strconv"
)

type Pager struct {
	Count   int
	Max     int
	Current int
	Enabled bool
	IsFirst bool
	IsLast  bool
	Prev    string
	Next    string
}

func pagerInit(itemCount, pageSize, pageNum int, prefix string) Pager {
	pages := 0
	if pageSize > 0 {
		pages = int(itemCount / pageSize)
		if itemCount%pageSize > 0 {
			pages++
		}
	}
	max := pages - 1
	prev := "#"
	if pageNum == 1 {
		prev = "../" + prefix + "/"
	} else if pageNum > 1 {
		prev = strconv.Itoa(pageNum-1) + ".html"
	}
	next := "#"
	if pageNum < max {
		next = strconv.Itoa(pageNum+1) + ".html"
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
	return p
}
