package main

import (
	// "flag"
	model "greenlight_gbolahan/internal/data"
	"greenlight_gbolahan/internal/jsonlog"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const version = "2.0.0"

type db struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type config struct {
	port        int
	env         string
	db          db
	rateLimiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models model.Models
}

func main() {

	cfg := startFlag()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)


	if err != nil {
		logger.PrintFatal(err, nil)
		log.Print("error is here")
	}

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: model.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}
