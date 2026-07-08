package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type CoinGeckoFetcher struct {
	tickerDict map[string]string
}

func NewCoinGeckoFetcher() CoinGeckoFetcher {
	return CoinGeckoFetcher{
		tickerDict: map[string]string{
			"BTCUSDT": "bitcoin",
			"ETHUSDT": "ethereum",
		},
	}
}

func (c CoinGeckoFetcher) FetchPrice(ctx context.Context, symbol string) (float64, error) {
	client := &http.Client{}

	id, ok := c.tickerDict[symbol]
	if !ok {
		log.Printf("FetchPrice: incorrect symbol\n")
		return 0, fmt.Errorf("incorrect symbol")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", id), nil)
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

	var info map[string]map[string]float64
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Printf("FetchPrice: %v", err)
		return 0, err
	}

	price := info[id]["usd"]
	return price, nil
}

func (c CoinGeckoFetcher) Name() string {
	return "coingecko"
}
