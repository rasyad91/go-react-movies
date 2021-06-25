package main

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/repository/postgres"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var (
	cfg          config.Application
	idleTimeout  = time.Minute
	readTimeout  = 5 * time.Second
	writeTimeout = 5 * time.Second
	dbTimeout    = 5 * time.Second
)

func main() {

	cfg = config.Application{}

	flag.IntVar(&cfg.Port, "port", 4000, "backend server port to listen on")
	flag.StringVar(&cfg.Env, "env", "development", "Application environment (development|production")
	flag.StringVar(&cfg.Version, "version", "1.0.0", "Application version")
	flag.StringVar(&cfg.DB.DSN, "db_dsn", "postgres://rasyad@localhost:5432/go_movies?sslmode=disable", "DB connection string")
	flag.StringVar(&cfg.DB.Driver, "db_driver", "postgres", "DB Driver (postgres|mysql etc)")
	flag.StringVar(&cfg.JWT.Secret, "jwt_secret", "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160", "jwt secret")

	flag.Parse()

	cfg.Logger = log.New(os.Stdout, "INFO:", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		cfg.Logger.Fatalln(err)
	}
	defer db.Close()

	r := postgres.New(db)
	h := handler.Initialize(&cfg, r)
	handler.New(h)

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

func openDB(cfg config.Application) (*sql.DB, error) {
	db, err := sql.Open(cfg.DB.Driver, cfg.DB.DSN)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
