package geo_test

import (
	"testing"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"
	"github.com/devprod-tech/webike_recomendations-Vitalya/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGeoService(t *testing.T) {
	testCases := []struct {
		name         string
		analyticsMsg models.AnalyticsMessage
		geoInput     models.GeoInput
		nearestTO    string
		sharedRecMsg *models.RecommendationsMessage
		successGet   bool
	}{
		{
			name: "success",
			analyticsMsg: models.AnalyticsMessage{
				TimeToCrash: "test",
				Payload: models.AnalyticsPayloadField{
					Location: models.AnalyticsLocationField{
						Latitude:  55.638248,
						Longitude: 37.612877,
					},
				},
			},
			geoInput:     models.GeoInput{Latitude: 55.638248, Longitude: 37.612877},
			nearestTO:    "Bicycle Repair Shop",
			sharedRecMsg: &models.RecommendationsMessage{},
			successGet:   true,
		},
		{
			name: "fail",
			analyticsMsg: models.AnalyticsMessage{
				TimeToCrash: "test",
				Payload: models.AnalyticsPayloadField{
					Location: models.AnalyticsLocationField{
						Latitude:  55.638248,
						Longitude: 37.612877,
					},
				},
			},
			geoInput:     models.GeoInput{Latitude: 55.638248, Longitude: 37.612877},
			sharedRecMsg: &models.RecommendationsMessage{},
			successGet:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			mockGeoPort := mocks.NewGeoPort(t)
			mockGeoPort.
				On("Accept", tc.analyticsMsg).
				Return(tc.geoInput).
				Once()

			mockGeoAPIAdapter := mocks.NewGeoAPIAdapter(t)
			if !tc.successGet {
				mockGeoAPIAdapter.
					On("GetNearestTO", tc.geoInput.Longitude, tc.geoInput.Latitude).
					Return(nil).
					Once()
			} else {
				mockGeoAPIAdapter.
					On("GetNearestTO", tc.geoInput.Longitude, tc.geoInput.Latitude).
					Return(&models.GeoAPIResponse{
						DisplayName: tc.nearestTO,
					}).
					Once()
			}
			mockGeoExternalAdapter := mocks.NewGeoExternalAdapter(t)
			if tc.successGet {
				mockGeoExternalAdapter.
					On("SendToShared", models.GeoOutput{Nearest_TO: tc.nearestTO}, tc.sharedRecMsg).
					Once()
			}

			geoInput := mockGeoPort.Accept(tc.analyticsMsg)
			assert.Equal(t, geoInput, tc.geoInput)

			data := mockGeoAPIAdapter.GetNearestTO(geoInput.Longitude, geoInput.Latitude)
			if tc.successGet {
				assert.NotNil(t, data)
				assert.Equal(t, data.DisplayName, tc.nearestTO)
				geoOutput := models.GeoOutput{Nearest_TO: tc.nearestTO}
				mockGeoExternalAdapter.SendToShared(geoOutput, tc.sharedRecMsg)
			} else {
				assert.Nil(t, data)
			}
		})
	}
}
