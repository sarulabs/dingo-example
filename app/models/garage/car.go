package garage

// Car is the structure representing a car.
type Car struct {
	ID    string `json:"id" bson:"id"`
	Brand string `json:"brand" json:"brand"`
	Color string `json:"color" json:"color"`
}
