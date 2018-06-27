package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sarulabs/dingo-example/app/handlers"
	"github.com/sarulabs/dingo-example/app/middlewares"
	"github.com/sarulabs/dingo-example/var/lib/services/dic"
)

func main() {
	// create the app container
	app, err := dic.NewContainer()
	if err != nil {
		log.Fatal(err)
	}
	defer app.Delete()

	// create the http server
	r := mux.NewRouter()

	// middleware to add the container in the http requests
	m := middlewares.DingoMiddleware

	r.HandleFunc("/cars", m(handlers.GetCarListHandler, app)).Methods("GET")
	r.HandleFunc("/cars", m(handlers.PostCarHandler, app)).Methods("POST")
	r.HandleFunc("/cars/{carId}", m(handlers.GetCarHandler, app)).Methods("GET")
	r.HandleFunc("/cars/{carId}", m(handlers.PutCarHandler, app)).Methods("PUT")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:" + os.Getenv("SERVER_PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
