package handler

import (
	"crypto-aggregator/internal/aggregator"
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	agg *aggregator.Aggregator
}

func NewHandler(agg *aggregator.Aggregator) *Handler {
	return &Handler{agg: agg}
}

type PriceResponse struct {
	Average   float64            `json:"average"`
	Sources   map[string]float64 `json:"sources"`
	UpdatedAt int64              `json:"updated_at"`
}

func (h Handler) Prices(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	prices := make(map[string]PriceResponse)
	data := h.agg.GetData()

	for symbol, info := range data {
		prices[symbol] = PriceResponse{
			Average:   info.Average,
			Sources:   info.Sources,
			UpdatedAt: info.UpdatedAt.UnixMilli(),
		}
	}

	toReturn, err := json.Marshal(prices)
	if err != nil {
		log.Printf("[handler.Prices]: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Write(toReturn)
}
