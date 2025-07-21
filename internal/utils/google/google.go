package google

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"
)

const GooglePlacesURL = "https://places.googleapis.com/v1/places:searchNearby"

func CreateGoogleAPIRequest(ctx context.Context, googleAPIKey string, longitude, latitude float64) (*http.Request, error) {
	payload := createGooglePlacesPayload(longitude, latitude)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", GooglePlacesURL, bytes.NewBuffer(payloadJSON))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", googleAPIKey)
	req.Header.Set("X-Goog-FieldMask", "places.displayName")
	return req, nil
}

func createGooglePlacesPayload(longitude float64, latitude float64) models.GooglePlacesPayload {
	return models.GooglePlacesPayload{
		IncludedTypes:  []string{"car_repair"},
		MaxResultCount: 1,
		LocationRestriction: models.GoogleLocationRestriction{
			Circle: models.GoogleLocationCircle{
				Center: models.GoogleLocationCenter{
					Latitude:  latitude,
					Longitude: longitude,
				},
				Radius: 5000.0,
			},
		},
		RankPreference: "DISTANCE",
	}
}

func CreateGeoAPIResponseFromGoogleAPIResponse(respData models.GooglePlacesResponse) models.GeoAPIResponse {
	data := models.GeoAPIResponse{
		DisplayName: respData.Places[0].DisplayName.Text,
		Address:     respData.Places[0].Address,
		PhoneNumber: respData.Places[0].NationalPhoneNumber,
	}

	if data.PhoneNumber == "" {
		data.PhoneNumber = respData.Places[0].InternationalPhoneNumber
	}

	return data
}
