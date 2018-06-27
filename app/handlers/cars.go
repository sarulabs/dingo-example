package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sarulabs/dingo-example/app/models/helpers"
	"github.com/sarulabs/dingo-example/var/lib/services/dic"
)

// GetCarListHandler is the handler that prints the list of cars.
func GetCarListHandler(w http.ResponseWriter, r *http.Request) {
	cars, err := dic.CarRepository(r).FindAll()
	if err != nil {
		dic.Logger(r).Error(err.Error())
		helpers.JSONResponse(w, 500, "Internal Error")
		return
	}

	helpers.JSONResponse(w, 200, cars)
}

// PostCarHandler is the handler that adds a new car.
func PostCarHandler(w http.ResponseWriter, r *http.Request) {
	var data *Car
	helpers.ReadJSONBody(r, &data)

	car, err := dic.CarManager(r).Create(data)
	if err != nil {
		helpers.JSONResponse(w, err.Status(), map[string]interface{}{
			"error": err.PublicMessage(),
		})
		return
	}

	helpers.JSONResponse(w, 200, car)
}

// GetCarHandler is the handler that prints the characteristics of a car.
func GetCarHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["carId"]

	car, err := dic.CarManager(r).Get(id)
	if err != nil {
		helpers.JSONResponse(w, err.Status(), map[string]interface{}{
			"error": err.PublicMessage(),
		})
		return
	}

	helpers.JSONResponse(w, 200, car)
}

// PutCarHandler is the handler that updates a car.
func PutCarHandler(w http.ResponseWriter, r *http.Request) {
	var data *Car
	helpers.ReadJSONBody(r, &data)

	id := mux.Vars(r)["carId"]

	car, err := dic.CarManager(r).Update(id, data)
	if err != nil {
		helpers.JSONResponse(w, err.Status(), map[string]interface{}{
			"error": err.PublicMessage(),
		})
		return
	}

	helpers.JSONResponse(w, 200, car)
}
