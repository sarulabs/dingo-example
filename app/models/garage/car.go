package garage

// Car is the structure representing a car.
type Car struct {
	ID    string `json:"id,omitempty"`
	Brand string `json:"brand,omitempty"`
	Color string `json:"color,omitempty"`
}
