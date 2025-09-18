// Filename: cmd/api/main.go
// page 247

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Lee26Ed/qod/internal/data"
	_ "github.com/lib/pq"
)

const version = "1.0.0"
type configuration struct {
	port int 
	env string
	db struct {
		dsn string
	}
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config configuration
	logger *slog.Logger
	quoteModel data.QuoteModel
}



func loadConfig() configuration {
	var cfg configuration

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")
	// read in the dsn
    flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://quotes:fishsticks@localhost/quotes",
                  "PostgreSQL DSN")
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)",
              func(val string) error {
                   cfg.cors.trustedOrigins = strings.Fields(val)
                   return nil
              })

	flag.Parse()

	return cfg
}

func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return logger
}

func openDB(settings configuration) (*sql.DB, error) {
    // open a connection pool
    db, err := sql.Open("postgres", settings.db.dsn)
    if err != nil {
        return nil, err
    }
    
    // set a context to ensure DB operations don't take too long
    ctx, cancel := context.WithTimeout(context.Background(),
                                       5 * time.Second)
    defer cancel()

    // let's test if the connection pool was created
    // we trying pinging it with a 5-second timeout
    err = db.PingContext(ctx)
    if err != nil {
        db.Close()
        return nil, err
    }

    // return the connection pool (sql.DB)
    return db, nil
} 

func main() {
	// Initialize configuration
	cfg := loadConfig()

	// Initialize logger
	logger := setupLogger()

	// the call to openDB() sets up our connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// release the database resources before exiting
	defer db.Close()

	logger.Info("database connection pool established")

	// Initialize application with dependencies
	app := &application {
		config: cfg,
		logger: logger,
		quoteModel: data.QuoteModel{DB: db},
	}

	apiServer := &http.Server {
        Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
        IdleTimeout: time.Minute,
        ReadTimeout: 5 * time.Second,
        WriteTimeout: 10 * time.Second,
        ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
    }

	logger.Info("starting server", "address", apiServer.Addr,
                "environment", cfg.env)
    err = apiServer.ListenAndServe()   // remove the :
    logger.Error(err.Error())
    os.Exit(1)
}

func printUB() string {
    return "Hello, UB!"
}