package geo

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/adapters"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/observability"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/utils/google"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

type GeoGoogleAdapter struct {
	log                 *zap.Logger
	client              *http.Client
	workStatusCommitter adapters.WorkStatusCommitter
	googleAPIKey        string
}

func NewGoogle(logger *zap.Logger, client *http.Client, googleAPIKey string) *GeoGoogleAdapter {
	logger = logger.With(zap.String("adapter", "geo"))
	if client == nil {
		client = &http.Client{Timeout: time.Second * 5}
	}
	client.Transport = otelhttp.NewTransport(http.DefaultTransport)
	return &GeoGoogleAdapter{
		log:          logger,
		client:       client,
		googleAPIKey: googleAPIKey,
	}
}

func (g *GeoGoogleAdapter) SetupWorkStatusCommitter(workStatusCommitter adapters.WorkStatusCommitter) {
	g.workStatusCommitter = workStatusCommitter
}

func (g *GeoGoogleAdapter) GetNearestTO(ctx context.Context, longitude float64, latitude float64) (*models.GeoAPIResponse, error) {
	var err error

	if g.workStatusCommitter != nil {
		defer func() {
			if err != nil {
				g.workStatusCommitter.CommitWorkStatus(observability.StatusNearestTOError)
			} else {
				g.workStatusCommitter.CommitWorkStatus(observability.StatusNearestTOFound)
			}
		}()
	}

	req, err := google.CreateGoogleAPIRequest(ctx, g.googleAPIKey, longitude, latitude)
	if err != nil {
		g.log.Error(google.ErrCreateGoogleAPIRequest.Error(), zap.Error(err))
		return nil, google.ErrCreateGoogleAPIRequest
	}

	resp, err := g.client.Do(req)
	if err != nil {
		g.log.Error(google.ErrGetNearestTOViaGoogleAPI.Error(), zap.Error(err))
		return nil, google.ErrGetNearestTOViaGoogleAPI
	}
	if resp.StatusCode != http.StatusOK {
		g.log.Error(google.ErrGetNearestTOViaGoogleAPI.Error(), zap.Int("status code", resp.StatusCode))
		return nil, google.ErrGetNearestTOViaGoogleAPI
	}

	defer resp.Body.Close()

	var respData models.GooglePlacesResponse

	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		g.log.Error(google.ErrDecodeResponseFromGoogleAPI.Error(), zap.Error(err))
		return nil, google.ErrDecodeResponseFromGoogleAPI
	}

	if len(respData.Places) == 0 {
		g.log.Info(google.ErrNoTOFoundViaGoogleAPI.Error(), zap.Float64("latitude", latitude), zap.Float64("longitude", longitude))
		return nil, google.ErrNoTOFoundViaGoogleAPI
	}

	data := google.CreateGeoAPIResponseFromGoogleAPIResponse(respData)
	return &data, nil
}
