package coingecko

import (
	"net/http"
	"time"

	"code.vegaprotocol.io/vega/logging"
	"github.com/vegaprotocol/vega-monitoring/config"
	"golang.org/x/time/rate"
)

type CoingeckoClient struct {
	httpClient  *http.Client
	config      *config.CoingeckoConfig
	rateLimiter *rate.Limiter
	log         *logging.Logger
}

func NewCoingeckoClient(config *config.CoingeckoConfig, log *logging.Logger) *CoingeckoClient {
	return &CoingeckoClient{
		config:      config,
		rateLimiter: rate.NewLimiter(rate.Every(200*time.Millisecond), 1),
		httpClient: &http.Client{
			Timeout: 2 * time.Second,
		},
		log: log,
	}
}
