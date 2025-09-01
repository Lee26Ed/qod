// Filename: cmd/api/main.go

package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"
type configuration struct {
	port int 
	env string
}

type application struct {
	config configuration
	logger *slog.Logger
}

func loadConfig() configuration {
	var cfg configuration

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")
	flag.Parse()

	return cfg
}

func setupLogger() *slog.Logger {
	var logger *slog.Logger 
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	return logger
}

func main() {
	// Initialize configuration
	cfg := loadConfig()

	// Initialize logger
	logger := setupLogger()

	// Initialize application with dependencies
	app := &application {
		config: cfg,
		logger: logger,
	}

	apiServer := &http.Server {
        Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
        IdleTimeout: time.Minute,
        ReadTimeout: 5 * time.Second,
        WriteTimeout: 10 * time.Second,
        ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
    }


	// Run the app
	err := apiServer.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func printUB() string {
    return "Hello, UB!"
}