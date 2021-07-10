package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg, err := loadConfig()
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
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server startup failed")
	}
}
func Greet(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello World!")
}
