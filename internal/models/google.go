package models

type GooglePlacesPayload struct {
	IncludedTypes       []string                  `json:"includedTypes"`
	MaxResultCount      int                       `json:"maxResultCount"`
	LocationRestriction GoogleLocationRestriction `json:"locationRestriction"`
	RankPreference      string                    `json:"rankPreference"`
}

type GoogleLocationRestriction struct {
	Circle GoogleLocationCircle `json:"circle"`
}

type GoogleLocationCircle struct {
	Center GoogleLocationCenter `json:"center"`
	Radius float64              `json:"radius"`
}

type GoogleLocationCenter struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type GooglePlacesResponse struct {
	Places []GooglePlacesResult `json:"places"`
}

type GooglePlacesResult struct {
	DisplayName              GoogleDisplayName `json:"displayName"`
	Address                  string            `json:"shortFormattedAddress"`
	NationalPhoneNumber      string            `json:"nationalPhoneNumber"`
	InternationalPhoneNumber string            `json:"internationalPhoneNumber"`
}

type GoogleDisplayName struct {
	Text string `json:"text"`
}
