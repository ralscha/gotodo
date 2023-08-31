package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         app.config.HTTP.Port,
		Handler:      app.routes(),
		ReadTimeout:  time.Duration(app.config.HTTP.ReadTimeoutInSeconds) * time.Second,
		WriteTimeout: time.Duration(app.config.HTTP.WriteTimeoutInSeconds) * time.Second,
		IdleTimeout:  time.Duration(app.config.HTTP.IdleTimeoutInSeconds) * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		shutdownChannel := app.taskScheduler.Shutdown()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		<-shutdownChannel
		app.wg.Wait()

		shutdownError <- nil
	}()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return <-shutdownError
}

func (app *application) background(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(error); ok {
					slog.Error("background job failed", e)
				} else {
					slog.Error("background job failed", nil, err)
				}
			}
		}()

		fn()
	}()
}
