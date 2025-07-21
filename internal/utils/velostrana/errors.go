package velostrana

import "errors"

var (
	ErrFailedToCreateURLToListDetails = errors.New("failed to create URL to list details")
	ErrFailedTORequestMarket          = errors.New("failed to request market")
	ErrFailedToParseMarket            = errors.New("failed to parse market")
	ErrCreateVelostranaRequest        = errors.New("failed to create velostrana request")
)
