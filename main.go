//+build !test

// This file holds code which does not need tests

// Package main contains only app.Run(os.Exit) call
// see https://github.com/apisite/apisite/app for details
package main

import (
	"log"
	"os"
)

// Actual version value will be set at build time
var version = "0.0-dev"

func main() {
	cfg, err := setupConfig()
	if err != nil {
		if err.Error() == "ERR1" {
			os.Exit(1)
		}
		os.Exit(2)
	}
	l := setupLog()
	r, err := setupRouter(cfg, l)
	if err != nil {
		log.Fatal(err)
	}
	err = runServer(l, cfg.Addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
