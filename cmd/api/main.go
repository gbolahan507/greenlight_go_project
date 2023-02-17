package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

)

// Declare a string containing the application version number.
const version = "2.0.0"

// Define a config struct to hold all the configuration settings for our application.
type config struct {
	port int
	env  string
}


type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API Server Port")
	flag.StringVar(&cfg.env, "env", "development", "Enviroment(development|staging|production)")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	route := app.routes()

	serve := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      route,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, serve.Addr)
	err := serve.ListenAndServe()
	log.Fatal(err)

}

// using flag

// go run ./cmd/api -port=3030 -env=production


// Read on

// json decoding nuances   pg 80
