package models

type GeoInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type GeoOutput struct {
	Nearest_TO string `json:"nearest_TO"`
}

type GeoAPIResponse struct {
	DisplayName string `json:"displayName"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phoneNumber"`
}
