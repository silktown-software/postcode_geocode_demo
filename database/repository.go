package database

import "errors"

// ErrEmptyPostcode a sentinel error representing that an empty string was provided to the Get method
var ErrEmptyPostcode = errors.New("the postcode is an empty string")

// PostCodeEntity this entity which holds the data for Postcode and Lng, Lat
type PostCodeEntity struct {
	Postcode string  `json:"postcode"`
	Lng      float64 `json:"lng"`
	Lat      float64 `json:"lat"`
}

// PostCodeRepository interface defining the methods which are required for the PostcodeRepository
type PostCodeRepository interface {
	Migrate() error
	Get(postcode string) (*PostCodeEntity, error)
	Upsert(entity PostCodeEntity) error
	UpsertMany(postcodes []PostCodeEntity) error
}
