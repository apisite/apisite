// +build !test

// This file holds code which does not covered by tests

/*
Package main contains only `app.Run(os.Exit)` call.

See https://godoc.org/github.com/apisite/apisite/app for details.

Command options:

  $ ./apisite -h
  apisite v0.6.1-6-ga26c508. API website engine

Usage:
  apisite [OPTIONS]

Application Options:
      --http_addr=      Http listen address (default: localhost:8080)
      --content_type=   Default content type (default: text/html; charset=utf-8)
      --error_404=      Template called when page is not found (default: .404)
      --buffer_size=    Template buffer size (default: 64)

Filesystem Options:
      --fs.templates=   Templates root path (default: tmpl/)
      --fs.mask=        Templates filename mask (default: .tmpl)
      --fs.includes=    Includes path (default: inc/)
      --fs.layouts=     Layouts path (default: layout/)
      --fs.pages=       Pages path (default: page/)
      --fs.use_suffix   Template type defined by suffix
      --fs.index=       Index page name (default: index)
      --fs.def_layout=  Default layout template (default: default)
      --fs.hide_prefix= Treat files with this prefix as hidden (default: .)

API Options:
      --api.dsn=        Database connect string (default: postgres://?sslmode=disable)
      --api.driver=     Database driver (default: postgres)
      --api.indef=      Argument definition function (default: func_args)
      --api.outdef=     Result row definition function (default: func_result)
      --api.index=      Available functions list (default: index)
      --api.schema=     Definition functions schema (default: rpc)
      --api.arg_syntax= Default named args syntax (:= or =>) (default: :=)
      --api.arg_prefix= Trim prefix from arg name (default: a_)
      --api.nsp=        Proc namespace(s)

Help Options:
  -h, --help            Show this help message

*/
package main

import (
	"github.com/apisite/apisite/app"
	"log"
	"os"
)

// Actual version value will be set at build time
var version = "0.0-dev"

func main() {
	app.Version = version
	log.Printf("apisite %s. API website engine", version)
	app.Run(os.Exit)
}
