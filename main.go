package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	mapper "github.com/birkirb/loggers-mapper-logrus"
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"

	"github.com/jessevdk/go-flags"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	mulate "github.com/apisite/gin-mulate"
	pgfc "github.com/apisite/gin-pgfc"
)

// Config holds all config vars
type Config struct {
	Addr      string        `long:"http_addr" default:"localhost:8080"  description:"Http listen address"`
	DBConnect string        `long:"db_connect" default:"" description:"Database connect string, i.e. user:pass@host/dbname?sslmode=disable"`
	Template  mulate.Config `group:"Template Options"`
	PGFC      pgfc.Config   `group:"PGFC Options"`
}

func main() {

	cfg := &Config{}
	p := flags.NewParser(cfg, flags.Default)

	_, err := p.Parse()
	if err != nil {
		if !strings.HasPrefix(err.Error(), "\nUsage") {
			fmt.Fprintf(os.Stderr, "error: %+v", err)
		}
		os.Exit(0)
	}

	l := logrus.New()

	if gin.IsDebugging() {
		l.SetLevel(logrus.DebugLevel)
		l.AddHook(filename.NewHook())
	}
	log := mapper.NewLogger(l)

	r := gin.Default()

	templates := mulate.New(cfg.Template, log)
	templates.DisableCache(gin.IsDebugging())

	allFuncs := template.FuncMap{}
	appendFuncs(allFuncs)

	s, err := pgfc.NewServer(cfg.PGFC, log, cfg.DBConnect, nil)
	if err != nil {
		log.Fatal(err)
	}
	s.SetFuncBlank(allFuncs)
	err = templates.LoadTemplates(allFuncs)
	if err != nil {
		log.Fatal(err)
	}

	templates.Route("", r)

	r.Use(static.Serve("/", static.LocalFile("./static", false)))
	r.NoRoute(func(c *gin.Context) {
		c.File("static/index.html")
	})

	s.Route("/rpc", r)

	templates.FuncHandler = func(ctx *gin.Context, funcs template.FuncMap) {
		s.SetFuncRequest(funcs, ctx)
	}

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server stopped")
}
