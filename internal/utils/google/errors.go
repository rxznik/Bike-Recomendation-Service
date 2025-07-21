package google

import "errors"

var (
	ErrCreateGoogleAPIRequest      = errors.New("failed to create google api request")
	ErrGetNearestTOViaGoogleAPI    = errors.New("failed to get nearest to via google api")
	ErrNoTOFoundViaGoogleAPI       = errors.New("no TO found via google api")
	ErrDecodeResponseFromGoogleAPI = errors.New("failed to decode response from google api")
)
