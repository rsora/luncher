package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ardanlabs/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var build = "dev"

func main() {

	// Construct application logger.
	log, err := initLogger("luncher")
	if err != nil {
		fmt.Println("Error constructing logger:", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Errorw("startup", "error", err)
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	// Define configuration elements and their defaults.
	cfg := struct {
		conf.Version
		Web struct {
			APIHost string `conf:"default:0.0.0.0:8000"`
		}
	}{
		Version: conf.Version{
			SVN:  build,
			Desc: "This service is supposed to provide the NST recipes.",
		},
	}

	// Configuration parsing
	const prefix = "luncher"
	help, err := conf.ParseOSArgs(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config error: %w", err)
	}

	// Configuration print on logs.
	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	return nil
}

func initLogger(service string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": service,
	}

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}
