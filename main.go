package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sarulabs/dingo-example/app/handlers"
	"github.com/sarulabs/dingo-example/app/middlewares"
	"github.com/sarulabs/dingo-example/config/logging"
	"github.com/sarulabs/dingo-example/var/lib/services/dic"
)

func main() {
	// Use a single logger in the whole application.
	// Need to close it at the end.
	defer logging.Logger.Sync()

	// Create the app container.
	app, err := dic.NewContainer()
	if err != nil {
		logging.Logger.Fatal(err.Error())
	}
	defer app.Delete()

	dic.ErrorCallback = func(err error) {
		logging.Logger.Error("Dingo error: " + err.Error())
	}

	// Create the http server.
	r := mux.NewRouter()

	// Function to apply the middlewares:
	// - recover from panic
	// - add the container in the http requests
	m := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.PanicRecoveryMiddleware(
			middlewares.DingoMiddleware(h, app),
			logging.Logger,
		)
	}

	r.HandleFunc("/cars", m(handlers.GetCarListHandler)).Methods("GET")
	r.HandleFunc("/cars", m(handlers.PostCarHandler)).Methods("POST")
	r.HandleFunc("/cars/{carId}", m(handlers.GetCarHandler)).Methods("GET")
	r.HandleFunc("/cars/{carId}", m(handlers.PutCarHandler)).Methods("PUT")
	r.HandleFunc("/cars/{carId}", m(handlers.DeleteCarHandler)).Methods("DELETE")

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:" + os.Getenv("SERVER_PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logging.Logger.Info("listening on port " + os.Getenv("SERVER_PORT"))

	if err := srv.ListenAndServe(); err != nil {
		logging.Logger.Error(err.Error())
	}
}
