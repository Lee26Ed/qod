// Filename: cmd/api/main.go

package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
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

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	return router
}

func (app *application) serve() error {
	srv := &http.Server {
		Addr: fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Minute,
		WriteTimeout: 10 * time.Second,
		ErrorLog: slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}
	app.logger.Info("Starting Server", "addr", srv.Addr, "env", app.config.env)

	return srv.ListenAndServe()
}

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	js := `{"status": "available", "environment": %q, "version": %q}`
	js = fmt.Sprintf(js, app.config.env, version)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(js))
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

	// Run the app
	err := app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func printUB() string {
    return "Hello, UB!"
}