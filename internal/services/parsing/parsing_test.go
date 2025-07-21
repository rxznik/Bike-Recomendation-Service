package parsing_test

import (
	"testing"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"
	"github.com/devprod-tech/webike_recomendations-Vitalya/mocks"
	"github.com/stretchr/testify/assert"
)

func TestParsingService(t *testing.T) {
	testCases := []struct {
		name         string
		analyticsMsg models.AnalyticsMessage
		parsingInput models.ParsingInput
		market       string
		sharedRecMsg *models.RecommendationsMessage
		successGet   bool
	}{
		{
			name: "success",
			analyticsMsg: models.AnalyticsMessage{
				Detail:      "test",
				TimeToCrash: "test",
			},
			parsingInput: models.ParsingInput{Detail: "test"},
			market:       "https://example.com",
			sharedRecMsg: &models.RecommendationsMessage{},
			successGet:   true,
		},
		{
			name: "fail",
			analyticsMsg: models.AnalyticsMessage{
				Detail:      "test",
				TimeToCrash: "test",
			},
			parsingInput: models.ParsingInput{Detail: "test"},
			sharedRecMsg: &models.RecommendationsMessage{},
			successGet:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			mockParsingPort := mocks.NewParsingPort(t)
			mockParsingPort.
				On("Accept", tc.analyticsMsg).
				Return(tc.parsingInput).
				Once()

			mockParsingParseAdapter := mocks.NewParsingParseAdapter(t)
			if !tc.successGet {
				mockParsingParseAdapter.
					On("GetRelevantProduct", tc.parsingInput.Detail).
					Return(nil).
					Once()
			} else {
				mockParsingParseAdapter.
					On("GetRelevantProduct", tc.parsingInput.Detail).
					Return(&models.ParsingParseResponse{URL: tc.market}).
					Once()
			}

			mockParsingExternalAdapter := mocks.NewParsingExternalAdapter(t)
			if tc.successGet {
				mockParsingExternalAdapter.
					On("SendToShared", models.ParsingOutput{Market: tc.market}, tc.sharedRecMsg).
					Once()
			}

			parsingInput := mockParsingPort.Accept(tc.analyticsMsg)
			assert.Equal(t, tc.parsingInput, parsingInput)

			relevantProductFromMarket := mockParsingParseAdapter.GetRelevantProduct(parsingInput.Detail)
			if tc.successGet {
				assert.NotNil(t, relevantProductFromMarket)
				assert.Equal(t, relevantProductFromMarket.URL, tc.market)
				parsingOutput := models.ParsingOutput{Market: relevantProductFromMarket.URL}
				mockParsingExternalAdapter.SendToShared(parsingOutput, tc.sharedRecMsg)
			} else {
				assert.Nil(t, relevantProductFromMarket)
			}
		})
	}
}
