// Package app implements apisite command functions
package app

import (
	"context"
	"github.com/pkg/errors"
	"html/template"
	"log"
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

	"github.com/apisite/apisite/tplfunc"
	"github.com/apisite/apitpl"
	"github.com/apisite/apitpl/ginapitpl"
	"github.com/apisite/apitpl/lookupfs"
)

// Config holds all config vars
type Config struct {
	Addr           string        `long:"http_addr" default:"localhost:8080"  description:"Http listen address"`
	ContentType    string        `long:"content_type" default:"text/html; charset=utf-8" description:"Default content type"`
	Error404       string        `long:"error_404" default:".404" description:"Template called when page is not found"`
	BufferSize     int           `long:"buffer_size" default:"64" description:"Template buffer size"`
	MaxHeaderBytes int           `long:"maxheader" description:"MaxHeaderBytes"`
	ReadTimeout    time.Duration `long:"rto" default:"10s" description:"HTTP read timeout"`
	WriteTimeout   time.Duration `long:"wto" default:"60s" description:"HTTP write timeout"`

	FS  lookupfs.Config `group:"Filesystem Options" namespace:"fs" env-namespace:"FS"`
	API procapi.Config  `group:"API Options" namespace:"api" env-namespace:"API"`
}

var (
	// ErrGotHelp returned after showing requested help
	ErrGotHelp = errors.New("help printed")
	// ErrBadArgs returned after showing command args error message
	ErrBadArgs = errors.New("option error printed")
)

// Version value will be set in main()
var Version = "0.0-dev"

// Run called by main() and prerforms all of app
func Run(exitFunc func(code int)) {
	var err error
	var cfg *Config
	defer func() { shutdown(exitFunc, err) }()
	cfg, err = setupConfig()
	if err != nil {
		return
	}
	l := setupLog()
	var r *gin.Engine
	r, err = setupRouter(cfg, l)
	if err != nil {
		return
	}
	err = runServer(l, cfg, r)
}

// exit after deferred cleanups have run
func shutdown(exitFunc func(code int), e error) {
	if e != nil {
		var code int
		switch e {
		case ErrGotHelp:
			code = 3
		case ErrBadArgs:
			code = 2
		default:
			code = 1
			log.Printf("Run error: %s", e.Error())
		}
		exitFunc(code)
	}
}

func setupConfig() (*Config, error) {
	cfg := &Config{}
	p := flags.NewParser(cfg, flags.Default)
	if _, err := p.Parse(); err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			return nil, ErrGotHelp
		}
		return nil, ErrBadArgs
	}
	return cfg, nil
}

// setupLog creates logger
func setupLog() loggers.Contextual {
	l := logrus.New()
	if gin.IsDebugging() {
		l.SetLevel(logrus.DebugLevel)
		l.SetReportCaller(true)
	}
	return &mapper.Logger{Logger: l} // Same as mapper.NewLogger(l) but without info log message
}

func setupRouter(cfg *Config, log loggers.Contextual) (*gin.Engine, error) {

	// procapi

	db := procapi.New(cfg.API, log, nil)
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

	allFuncs := template.FuncMap{
		"version": func() string { return Version },
	}
	tplfunc.SetSimpleFuncs(allFuncs)
	tplfunc.SetProtoFuncs(allFuncs)
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
		tplfunc.SetRequestFuncs(funcs, ctx)
		api.SetRequestFuncs(funcs, ctx)
		m := tplfunc.NewMeta(http.StatusOK, cfg.ContentType)
		m.SetLayout(cfg.FS.DefLayout)
		return m
	}

	gintpl.Route("", r)

	r.Use(static.Serve("/", static.LocalFile("./static", false)))
	r.NoRoute(func(c *gin.Context) {
		gintpl.HTML(c, cfg.Error404)
	})

	if err = api.Route("/rpc", r); err != nil {
		return nil, err
	}

	return r, nil
}

func runServer(log loggers.Contextual, cfg *Config, r *gin.Engine) error {

	srv := &http.Server{
		Addr:           cfg.Addr,
		Handler:        r,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
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
