package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sarulabs/dingo-example/app/models/garage"
	"github.com/sarulabs/dingo-example/app/models/helpers"
	"github.com/sarulabs/dingo-example/var/lib/services/dic"
)

// GetCarListHandler is the handler that prints the list of cars.
func GetCarListHandler(w http.ResponseWriter, r *http.Request) {
	cars, err := dic.CarRepository(r).FindAll()

	if err == nil {
		helpers.JSONResponse(w, 200, cars)
		return
	}

	helpers.JSONResponse(w, 500, map[string]interface{}{
		"error": "Internal Error",
	})
}

// PostCarHandler is the handler that adds a new car.
func PostCarHandler(w http.ResponseWriter, r *http.Request) {
	var input *Car

	helpers.ReadJSONBody(r, &input)

	car, err := dic.CarManager(r).Create(input)

	if err == nil {
		helpers.JSONResponse(w, 200, cars)
		return
	}

	switch e := err.(type) {
	case garage.ErrValidation:
		helpers.JSONResponse(w, 400, map[string]interface{}{
			"error": e.PublicMessage,
		})
	case garage.ErrNotFound:
		helpers.JSONResponse(w, 404, map[string]interface{}{
			"error": e.PublicMessage,
		})
	default:
		helpers.JSONResponse(w, 500, map[string]interface{}{
			"error": "Internal Error",
		})
	}
}

// GetCarHandler is the handler that prints the characteristics of a car.
func GetCarHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["carId"]

	car, err := dic.CarManager(r).Get(id)

	if err == nil {
		helpers.JSONResponse(w, 200, cars)
		return
	}

	switch e := err.(type) {
	case garage.ErrNotFound:
		helpers.JSONResponse(w, 404, map[string]interface{}{
			"error": e.PublicMessage,
		})
	default:
		helpers.JSONResponse(w, 500, map[string]interface{}{
			"error": "Internal Error",
		})
	}
}

// PutCarHandler is the handler that updates a car.
func PutCarHandler(w http.ResponseWriter, r *http.Request) {
	var input *Car

	helpers.ReadJSONBody(r, &input)

	id := mux.Vars(r)["carId"]

	car, err := dic.CarManager(r).Update(id, input)

	if err == nil {
		helpers.JSONResponse(w, 200, cars)
		return
	}

	switch e := err.(type) {
	case garage.ErrValidation:
		helpers.JSONResponse(w, 400, map[string]interface{}{
			"error": e.PublicMessage,
		})
	case garage.ErrNotFound:
		helpers.JSONResponse(w, 404, map[string]interface{}{
			"error": e.PublicMessage,
		})
	default:
		helpers.JSONResponse(w, 500, map[string]interface{}{
			"error": "Internal Error",
		})
	}
}

// DeleteCarHandler is the handler that removes a car from the database.
func DeleteCarHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["carId"]

	car, err := dic.CarManager(r).Delete(id)

	if err == nil {
		// TODO: just modify the response code
		helpers.JSONResponse(w, 204, "")
		return
	}

	switch e := err.(type) {
	case garage.ErrNotFound:
		helpers.JSONResponse(w, 404, map[string]interface{}{
			"error": e.PublicMessage,
		})
	default:
		helpers.JSONResponse(w, 500, map[string]interface{}{
			"error": "Internal Error",
		})
	}
}
