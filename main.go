//+build !test

package main

var version = "0.2-dev"

func main() {
	cfg := initConfig()
	log := initLog()
	r := initRouter(cfg, log)
	serve(log, cfg.Addr, r)
}
