package product

import "errors"

var (
	ErrNameIsEmpty        = errors.New("product name is empty.")
	ErrDescriptionIsEmpty = errors.New("product description is empty.")
	ErrIDIsEmpty          = errors.New("product id is empty.")
)
