package garage

import (
	"errors"
	"strings"

	"github.com/sarulabs/dingo-example/app/models/helpers"
)

// Car is the structure representing a car.
type Car struct {
	ID    string `json:"id" bson:"_id"`
	Brand string `json:"brand" bson:"brand"`
	Color string `json:"color" bson:"color"`
}

// colorsByBrand is the list of available options
// for the car brands and colors.
var colorsByBrand = map[string][]string{
	"audi":    []string{"black", "white", "yellow", "red"},
	"porsche": []string{"black", "white", "yellow", "red", "green"},
	"bmw":     []string{"black", "white", "blue"},
	"ferrari": []string{"red", "yellow"},
}

// brands returns the list of available brands.
func brands() []string {
	var brands []string
	for brand := range colorsByBrand {
		brands = append(brands, brand)
	}
	return brands
}

// ValidateCar checks if the car brand and color are available.
// If the car is not valid, an helpers.ErrValidation is returned.
func ValidateCar(car *Car) error {
	colors, ok := colorsByBrand[car.Brand]
	if !ok {
		msg := "Brand `" + car.Brand + "` does not exist. Available brands: " +
			strings.Join(brands(), ", ")

		return helpers.ErrValidation{
			Err:           errors.New(msg),
			PublicMessage: msg,
		}
	}

	for _, color := range colors {
		if color == car.Color {
			return nil
		}
	}

	msg := "Color `" + car.Color + "` does not exist for `" + car.Brand +
		"`. Available colors: " + strings.Join(colors, ", ")

	return helpers.ErrValidation{
		Err:           errors.New(msg),
		PublicMessage: msg,
	}
}
