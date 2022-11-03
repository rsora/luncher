package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/pkg/errors"
)

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
	log.Println("init")
	var cfg struct {
		Web struct {
			APIHost string `conf:default:0.0.0.0:8000`
		}
	}

	if usage, err := conf.Parse("LUNCHER", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}
	return nil
}
