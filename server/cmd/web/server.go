package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         app.config.Http.Port,
		Handler:      app.routes(),
		ReadTimeout:  time.Duration(app.config.Http.ReadTimeoutInSeconds) * time.Second,
		WriteTimeout: time.Duration(app.config.Http.WriteTimeoutInSeconds) * time.Second,
		IdleTimeout:  time.Duration(app.config.Http.IdleTimeoutInSeconds) * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Infof("caught signal: %s", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.Info("stopping scheduled jobs")
		app.scheduleStopChan <- struct{}{}

		app.logger.Info("completing background tasks")
		app.wg.Wait()

		shutdownError <- nil
	}()

	app.logger.Infow("starting server", "addr", srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Infow("server stopped", "addr", srv.Addr)

	_ = app.logger.Sync()

	return nil
}
