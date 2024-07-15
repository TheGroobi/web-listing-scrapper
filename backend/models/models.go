package models

type CarListing struct {
	Title    string
	Link     string
	Gearbox  string
	FuelType string
	Color    string
	Version  string
	Year     uint16
	Power    uint16
	Mileage  uint32
	ID       uint32
	Price    float64
}
