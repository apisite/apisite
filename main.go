//+build !test

// This file holds code which does not need tests

package main

// Actual version value is set during build
var version = "0.0-dev"

func main() {
	cfg := initConfig()
	log := initLog()
	r := initRouter(cfg, log)
	serve(log, cfg.Addr, r)
}
