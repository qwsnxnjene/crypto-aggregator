package handler

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) PricesWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[handler.PricesWS]: %v", err)
		return
	}
	defer conn.Close()

	for {
		select {
		case <-h.agg.Subscribe():
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
				log.Printf("[handler.PricesWS]: %v", err)
				return
			}

			conn.WriteMessage(websocket.TextMessage, toReturn)
		case <-r.Context().Done():
			return
		}
	}
}
