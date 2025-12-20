package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"stocky/internal/models"
	"stocky/internal/repository"
)

type PriceService struct {
	priceRepo *repository.StockPriceRepository
}

func NewPriceService(priceRepo *repository.StockPriceRepository) *PriceService {
	return &PriceService{priceRepo: priceRepo}
}

var mockPrices = map[string]decimal.Decimal{
	"RELIANCE": decimal.NewFromInt(2500),
	"TCS":      decimal.NewFromInt(3500),
	"INFY":     decimal.NewFromInt(1500),
	"HDFCBANK": decimal.NewFromInt(1700),
	"ICICIBANK": decimal.NewFromInt(950),
}

func (s *PriceService) FetchAndStorePrices(ctx context.Context) error {
	logrus.Info("Starting price fetch job")

	for symbol, basePrice := range mockPrices {
		variation := decimal.NewFromFloat(rand.Float64()*0.1 - 0.05)
		newPrice := basePrice.Mul(decimal.NewFromInt(1).Add(variation))

		newPrice = newPrice.Round(2)

		price := &models.StockPrice{
			Symbol:    symbol,
			Price:     newPrice,
			FetchedAt: time.Now(),
		}

		if err := s.priceRepo.Upsert(ctx, price); err != nil {
			logrus.WithError(err).WithField("symbol", symbol).Error("Failed to store price")
			continue
		}

		logrus.WithFields(logrus.Fields{
			"symbol": symbol,
			"price":  newPrice,
		}).Info("Price updated")
	}

	logrus.Info("Price fetch job completed")
	return nil
}

func (s *PriceService) GetLatestPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	price, err := s.priceRepo.GetLatest(ctx, symbol)
	if err != nil {
		return decimal.Zero, err
	}
	return price.Price, nil
}

