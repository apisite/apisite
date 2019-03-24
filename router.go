package main

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"time"

	mapper "github.com/birkirb/loggers-mapper-logrus"
	//	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
	"gopkg.in/birkirb/loggers.v1"

	"github.com/jessevdk/go-flags"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/apisite/pgcall"
	"github.com/apisite/pgcall/gin-pgcall"
	"github.com/apisite/pgcall/pgx-pgcall"

	"github.com/apisite/tpl2x"
	"github.com/apisite/tpl2x/gin-tpl2x"
	"github.com/apisite/tpl2x/lookupfs"
)

// Config holds all config vars
type Config struct {
	Addr        string `long:"http_addr" default:"localhost:8080"  description:"Http listen address"`
	DBConnect   string `long:"db_connect" default:"" description:"Database connect string, i.e. user:pass@host/dbname?sslmode=disable"`
	ContentType string `long:"content_type" default:"text/html; charset=utf-8" description:"Default content type"`
	BufferSize  int    `long:"buffer_size" default:"64" description:"Template buffer size"`

	FS  lookupfs.Config  `group:"Filesystem Options" namespace:"fs" env-namespace:"FS"`
	DB  pgxpgcall.Config `group:"DB Options" namespace:"db" env-namespace:"DB"`
	API pgcall.Config    `group:"API Options" namespace:"api" env-namespace:"API"`
}

func initConfig() *Config {
	cfg := &Config{}
	p := flags.NewParser(cfg, flags.Default)
	if _, err := p.Parse(); err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			os.Exit(1) // help printed
		} else {
			os.Exit(2) // error message written already
		}
	}
	return cfg
}

func initLog() loggers.Contextual {
	l := logrus.New()

	if gin.IsDebugging() {
		l.SetLevel(logrus.DebugLevel)
		l.SetReportCaller(true)
		//l.AddHook(filename.NewHook())
	}
	return mapper.NewLogger(l)
}

func initRouter(cfg *Config, log loggers.Contextual) *gin.Engine {

	// pgcall

	pg, err := pgxpgcall.New(cfg.DB, log)
	if err != nil {
		log.Fatal(err)
	}
	caller, err := pgcall.New(cfg.API, log, pg)
	if err != nil {
		log.Fatal(err)
	}

	s := ginpgcall.NewServer(log, caller)

	// tpl2x

	allFuncs := make(template.FuncMap, 0)
	initFuncs(allFuncs)
	protoFuncs(allFuncs)
	s.SetFuncBlank(allFuncs)

	//	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	fs := lookupfs.New(cfg.FS)
	// TODO:	templates.DisableCache(gin.IsDebugging())
	tfs, err := tpl2x.New(cfg.BufferSize).Funcs(allFuncs).LookupFS(fs).Parse()
	if err != nil {
		log.Fatal(err)
	}
	gintpl := gintpl2x.New(log, tfs)
	gintpl.RequestHandler = func(ctx *gin.Context, funcs template.FuncMap) gintpl2x.MetaData {
		requestFuncs(funcs, ctx)
		s.SetFuncRequest(funcs, ctx)
		return &Meta{status: http.StatusOK, contentType: cfg.ContentType, layout: cfg.FS.DefLayout}
	}

	gintpl.Route("", r)

	r.Use(static.Serve("/", static.LocalFile("./static", false)))
	r.NoRoute(func(c *gin.Context) {
		c.File("static/index.html")
	})

	s.Route("/rpc", r)

	return r
}

func serve(log loggers.Contextual, addr string, r *gin.Engine) {

	srv := &http.Server{
		Addr:    addr,
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
