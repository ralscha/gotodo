package main

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
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

		log.Info().Msgf("caught signal: %s", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		// log.Info().Msg("completing background tasks")

		// app.wg.Wait()
		shutdownError <- nil
	}()

	log.Info().Str("addr", srv.Addr).Msg("starting server")

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Info().Str("addr", srv.Addr).Msg("stopped server")

	return nil
}
