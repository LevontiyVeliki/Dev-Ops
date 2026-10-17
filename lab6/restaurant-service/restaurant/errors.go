package restaurant

import "errors"

var (
	ErrRestaurantNotFound = errors.New("restaurant not found")
	ErrInvalidName        = errors.New("restaurant name is required")
	ErrInvalidRating      = errors.New("rating must be between 0 and 5")
	ErrCmdRepository      = errors.New("unable to command repository")
)
