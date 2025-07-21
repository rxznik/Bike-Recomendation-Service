package parsing

import (
	"context"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/adapters"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/models"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/observability"
	"github.com/devprod-tech/webike_recomendations-Vitalya/internal/utils/velostrana"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

type ParsingParseAdapter struct {
	client              *http.Client
	log                 *zap.Logger
	workStatusCommitter adapters.WorkStatusCommitter
}

func NewParse(logger *zap.Logger, client *http.Client) *ParsingParseAdapter {
	logger = logger.With(zap.String("adapter", "parsing"))
	if client == nil {
		client = &http.Client{Timeout: time.Second * 5}
	}
	client.Transport = otelhttp.NewTransport(http.DefaultTransport)
	return &ParsingParseAdapter{
		log:    logger,
		client: client,
	}
}

func (p *ParsingParseAdapter) SetupWorkStatusCommitter(workStatusCommitter adapters.WorkStatusCommitter) {
	p.workStatusCommitter = workStatusCommitter
}

func (p *ParsingParseAdapter) GetRelevantProduct(ctx context.Context, detail string) (*models.ParsingParseResponse, error) {
	var err error

	if p.workStatusCommitter != nil {
		defer func() {
			if err != nil {
				p.workStatusCommitter.CommitWorkStatus(observability.StatusMarketError)
			} else {
				p.workStatusCommitter.CommitWorkStatus(observability.StatusMarketParsed)
			}
		}()
	}

	req, err := velostrana.CreateVelostranaRequest(ctx, detail)
	if err != nil {
		p.log.Error(err.Error(), zap.Error(err))
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		p.log.Error(velostrana.ErrFailedTORequestMarket.Error(), zap.Error(err))
		return nil, velostrana.ErrFailedTORequestMarket
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		p.log.Error(velostrana.ErrFailedTORequestMarket.Error(), zap.Int("status code", resp.StatusCode))
		return nil, velostrana.ErrFailedTORequestMarket
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		p.log.Error(velostrana.ErrFailedToParseMarket.Error(), zap.Error(err))
		return nil, velostrana.ErrFailedToParseMarket
	}
	firstProductHTMLCard := doc.Find(".product-card").First()

	productPath := firstProductHTMLCard.Find("a").First().AttrOr("href", "")
	p.log.Info("market parsed", zap.String("product path", productPath))

	return &models.ParsingParseResponse{
		URL: velostrana.AddBaseURL(productPath),
	}, nil
}
