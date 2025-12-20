package scheduler

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"stocky/internal/service"
)

type PriceFetcher struct {
	priceService *service.PriceService
	interval     time.Duration
}

func NewPriceFetcher(priceService *service.PriceService, interval time.Duration) *PriceFetcher {
	return &PriceFetcher{
		priceService: priceService,
		interval:     interval,
	}
}

func (pf *PriceFetcher) Start(ctx context.Context) {
	ticker := time.NewTicker(pf.interval)
	defer ticker.Stop()

	if err := pf.priceService.FetchAndStorePrices(ctx); err != nil {
		logrus.WithError(err).Error("Initial price fetch failed")
	}

	for {
		select {
		case <-ctx.Done():
			logrus.Info("Price fetcher stopped")
			return
		case <-ticker.C:
			if err := pf.priceService.FetchAndStorePrices(ctx); err != nil {
				logrus.WithError(err).Error("Price fetch failed")
			}
		}
	}
}

