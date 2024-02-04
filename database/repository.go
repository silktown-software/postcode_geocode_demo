package database

import "errors"

var ErrEmptyPostcode = errors.New("postcode is an empty string")

type PostCodeEntity struct {
	Postcode string  `json:"postcode"`
	Lng      float64 `json:"lng"`
	Lat      float64 `json:"lat"`
}

type PostCodeRepository interface {
	Migrate() error
	Get(postcode string) (*PostCodeEntity, error)
	Insert(entity PostCodeEntity) error
	InsertMany(postcodes []PostCodeEntity) error
}
