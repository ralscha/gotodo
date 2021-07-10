package main

import (
	"fmt"
	"gotodo.rasc.ch/internal/config"
	"log"
	"net/http"
	"time"
)

var (
	buildTime string
	version   string
)

func main() {
	fmt.Println(buildTime, version)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("reading config failed", err)
	}
	fmt.Println(cfg)
	mux := http.NewServeMux()
	mux.HandleFunc("/", Greet)
	log.Println("Starting server " + cfg.Http.Port)
	s := &http.Server{
		Addr:         cfg.Http.Port,
		Handler:      mux,
		ReadTimeout:  cfg.Http.ReadTimeoutInSeconds * time.Second,
		WriteTimeout: cfg.Http.WriteTimeoutInSeconds * time.Second,
		IdleTimeout:  cfg.Http.IdleTimeoutInSeconds * time.Second,
	}
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server startup failed")
	}
}
func Greet(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello World!")
}
