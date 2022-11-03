package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/ardanlabs/conf/v3"
	"github.com/pkg/errors"
)

// build is the git version of this program. It is set
// using build flags during the build process.
var build = "develop"

func main() {
	if err := run(); err != nil {
		log.Println("error: ", err)
		os.Exit(1)
	}

}

func run() error {
	log := log.New(os.Stdout, "LUNCHER : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// This is the config struct for the luncher-api app.
	var cfg struct {
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:8000"`
			DebugHost       string        `conf:"default:0.0.0.0:8001"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
	}

	// Using ardanlabs/conf package we get OOB also the `--help`
	// flag that prints all the config items required by the app.
	if usage, err := conf.Parse("LUNCHER", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// Print the build version for our logs, and
	// expose it under /debug/vars
	expvar.NewString("build").Set(build)
	log.Printf("main : Started : Application initializing : version %q", build)
	defer log.Println("main : Completed")

	// Print the config we are going to use.
	out, err := conf.String((&cfg))
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main : Config :\n%v\n", out)

	// Start a debug server.
	log.Println("main : debug Server : Starting")

	// Automagically, both expvar and pprof packages
	// add to the default MUX a /debug/pprof and a /debug/vars
	// routes.
	// Use `expvarmon` tool to inspect with pretty printing exported vars.
	go func() {
		log.Printf("main : Debug Listening %s", cfg.Web.DebugHost)
		log.Printf("main : Debug Listener closed", http.ListenAndServe(cfg.Web.DebugHost, http.DefaultServeMux()))
	}()

	// Start the API service
	log.Println("main : API Server : Starting")

	// Make a channel to listen	for interrupt or Terminate signal from OS.
	// Use a buffered channel because signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// a := App{}
	// a.Initialize()
	// a.Run(cfg.Web.APIHost)

	api := http.Server{
		Addr:    cfg.Web.APIHost,
		Handler: handlers.API(shutdown, log),
	}

	// Make a channel to listen to server errors.
	// Use a buffered channel so that the go routine can exit if
	// we don't fetch the error.
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("main : API Listening %s", api.Addr)
		serverErrors <- api.ListenAndServe()

	}()

	// Blocking main and wait for shutdown.
	select {

	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		log.Printf("main : %v : Shutdown started", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			err = api.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop gracefully")
		}
	}
	return nil
}
