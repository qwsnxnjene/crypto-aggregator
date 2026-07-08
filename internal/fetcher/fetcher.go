package fetcher

import (
	"context"
	"time"
)

type PriceFetcher interface {
	FetchPrice(ctx context.Context, symbol string) (float64, error)
	Name() string
}

type Price struct {
	Price     float64
	Symbol    string
	UpdatedAt time.Time
	Source    string
}
