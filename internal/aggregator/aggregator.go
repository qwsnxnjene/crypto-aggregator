package aggregator

import (
	"context"
	"crypto-aggregator/internal/fetcher"
	"crypto-aggregator/internal/storage"
	"log"
	"sync"
	"time"
)

type SymbolData struct {
	Sources   map[string]float64
	Average   float64
	UpdatedAt time.Time
}

type Aggregator struct {
	mu       sync.RWMutex
	fetchers []fetcher.PriceFetcher
	data     map[string]*SymbolData
	notify   chan struct{}
	notifyMu sync.RWMutex
	storage  *storage.Storage
}

func NewAggregator(storage *storage.Storage) Aggregator {
	return Aggregator{
		sync.RWMutex{},
		[]fetcher.PriceFetcher{fetcher.BinanceFetcher{}, fetcher.NewCoinGeckoFetcher()},
		map[string]*SymbolData{
			"BTCUSDT": {Sources: map[string]float64{}},
			"ETHUSDT": {Sources: map[string]float64{}},
		},
		make(chan struct{}),
		sync.RWMutex{},
		storage,
	}
}

func (a *Aggregator) Run(ctx context.Context) {
	for {
		if err := ctx.Err(); err != nil {
			return
		}

		priceChan := make(chan fetcher.Price)
		wg := sync.WaitGroup{}

		for symb := range a.data {
			for _, fetch := range a.fetchers {
				wg.Add(1)
				go func(symb string, fetch fetcher.PriceFetcher, priceChan chan fetcher.Price) {
					defer wg.Done()
					price, err := fetch.FetchPrice(ctx, symb)
					if err != nil {
						log.Printf("[aggregator.Run]: %v", err)
						return
					}
					priceChan <- fetcher.Price{
						Price:     price,
						Source:    fetch.Name(),
						UpdatedAt: time.Now(),
						Symbol:    symb,
					}
				}(symb, fetch, priceChan)
			}
		}

		go func() {
			wg.Wait()
			close(priceChan)
		}()

		for data := range priceChan {
			a.mu.Lock()
			a.data[data.Symbol].Sources[data.Source] = data.Price
			a.data[data.Symbol].UpdatedAt = data.UpdatedAt
			sum, count := 0.0, 0
			for _, p := range a.data[data.Symbol].Sources {
				sum += p
				count++
			}
			avg := sum / float64(count)
			a.data[data.Symbol].Average = avg
			// fmt.Printf("%s from %s at %v: %v\n", data.Symbol, data.Source, data.UpdatedAt, data.Price)
			a.mu.Unlock()
		}
		for symbol, info := range a.data {
			err := a.storage.Save(ctx, symbol, info.Average, info.Sources)
			if err != nil {
				log.Printf("[aggregator.Run]: can't save to DB: %v", err)
			}
		}
		a.notifyMu.Lock()
		close(a.notify)
		a.notify = make(chan struct{})
		a.notifyMu.Unlock()
		time.Sleep(time.Second * 5)
	}
}

func (a *Aggregator) GetData() map[string]SymbolData {
	defer a.mu.RUnlock()
	a.mu.RLock()

	res := make(map[string]SymbolData, len(a.data))
	for k, v := range a.data {
		res[k] = *v
	}

	return res
}

func (a *Aggregator) Subscribe() chan struct{} {
	a.notifyMu.RLock()
	defer a.notifyMu.RUnlock()
	return a.notify
}
