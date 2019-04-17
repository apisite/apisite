package main

import (
	"context"
	"github.com/pkg/errors"
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

	"github.com/apisite/procapi"
	"github.com/apisite/procapi/ginproc"

	"github.com/apisite/procapi/pgtype"
	_ "github.com/lib/pq"

	"github.com/apisite/apitpl"
	"github.com/apisite/apitpl/ginapitpl"
	"github.com/apisite/apitpl/lookupfs"
)

// Config holds all config vars
type Config struct {
	Addr        string `long:"http_addr" default:"localhost:8080"  description:"Http listen address"`
	ContentType string `long:"content_type" default:"text/html; charset=utf-8" description:"Default content type"`
	Error404    string `long:"error_404" default:".404" description:"Template called when page is not found"`
	BufferSize  int    `long:"buffer_size" default:"64" description:"Template buffer size"`

	FS  lookupfs.Config `group:"Filesystem Options" namespace:"fs" env-namespace:"FS"`
	API procapi.Config  `group:"API Options" namespace:"api" env-namespace:"API"`
}

func setupConfig() (*Config, error) {
	cfg := &Config{}
	p := flags.NewParser(cfg, flags.Default)
	if _, err := p.Parse(); err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			return nil, errors.New("ERR1") //os.Exit(1) // help printed
		} else {
			return nil, errors.New("ERR2") //os.Exit(2) // error message written already
		}
	}
	return cfg, nil
}

func setupLog() loggers.Contextual {
	l := logrus.New()

	if gin.IsDebugging() {
		l.SetLevel(logrus.DebugLevel)
		l.SetReportCaller(true)
		//l.AddHook(filename.NewHook())
	}
	return mapper.NewLogger(l)
}

func setupRouter(cfg *Config, log loggers.Contextual) (*gin.Engine, error) {

	// procapi

	db := procapi.New(cfg.API, log, nil).SetMarshaller(pgtype.New())
	err := db.Open()
	if err != nil {
		return nil, err
	}
	err = db.LoadMethods()
	if err != nil {
		return nil, err
	}
	api := ginproc.NewServer(log, db)

	// apitpl

	allFuncs := make(template.FuncMap, 0)
	SetSimpleFuncs(allFuncs)
	SetProtoFuncs(allFuncs)
	api.SetProtoFuncs(allFuncs)

	//	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	fs := lookupfs.New(cfg.FS)
	tfs, err := apitpl.New(cfg.BufferSize).Funcs(allFuncs).LookupFS(fs).Parse()
	if err != nil {
		return nil, err
	}
	tfs.ParseAlways(gin.IsDebugging())
	gintpl := ginapitpl.New(log, tfs)
	gintpl.RequestHandler = func(ctx *gin.Context, funcs template.FuncMap) ginapitpl.MetaData {
		SetRequestFuncs(funcs, ctx)
		api.SetRequestFuncs(funcs, ctx)
		m := NewMeta(http.StatusOK, cfg.ContentType)
		m.SetLayout(cfg.FS.DefLayout)
		return m
	}

	gintpl.Route("", r)

	r.Use(static.Serve("/", static.LocalFile("./static", false)))
	r.NoRoute(func(c *gin.Context) {
		gintpl.HTML(c, cfg.Error404)
	})

	api.Route("/rpc", r)

	return r, nil
}

func runServer(log loggers.Contextual, addr string, r *gin.Engine) error {

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
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
	return nil
}
