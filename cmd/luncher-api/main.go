package main

import (
	"expvar"
	"fmt"
	"log"
	"os"

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

	// a := App{}
	// a.Initialize()
	// log.Println("Starting...")
	// a.Run(":8000")
}

func run() error {
	log := log.New(os.Stdout, "LUNCHER : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// This is the config struct for the luncher-api app.
	var cfg struct {
		Web struct {
			APIHost string `conf:default:0.0.0.0:8000`
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

	return nil

}
