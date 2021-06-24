package main

import (
	"backend/internal/config"
	"backend/internal/handler"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	cfg          config.Application
	idleTimeout  = time.Minute
	readTimeout  = 5 * time.Second
	writeTimeout = 5 * time.Second
)

func main() {

	cfg = config.Application{}

	flag.IntVar(&cfg.Port, "port", 4000, "backend server port to listen on")
	flag.StringVar(&cfg.Env, "env", "development", "Application environment (development|production")
	flag.StringVar(&cfg.Version, "version", "1.0.0", "Application version")
	flag.Parse()

	cfg.Logger = log.New(os.Stdout, "INFO:", log.Ldate|log.Ltime)

	r := handler.Initialize(&cfg)
	handler.New(r)

	s := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      route(),
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	cfg.Logger.Println("Running on port ", cfg.Port)

	if err := s.ListenAndServe(); err != nil {
		cfg.Logger.Fatal(err)
	}

}
