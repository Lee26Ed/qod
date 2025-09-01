package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func (app *application) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}
	app.logger.Info("Starting Server", "addr", srv.Addr, "env", app.config.env)

	return srv.ListenAndServe()
}