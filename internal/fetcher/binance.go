package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type PriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type BinanceFetcher struct{}

func (b BinanceFetcher) FetchPrice(ctx context.Context, symbol string) (float64, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol), nil)
	if err != nil {
		log.Printf("FetchPrice: %v", err)
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("FetchPrice: %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("FetchPrice: %v", err)
		return 0, err
	}

	info := &PriceResponse{}
	err = json.Unmarshal(body, info)
	if err != nil {
		log.Printf("FetchPrice: %v", err)
		return 0, err
	}

	price, err := strconv.ParseFloat(info.Price, 64)
	if err != nil {
		log.Printf("FetchPrice: %v", err)
		return 0, err
	}

	return price, nil
}

func (b BinanceFetcher) Name() string {
	return "binance"
}
