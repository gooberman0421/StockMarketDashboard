package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type StockPrice struct {
	Ticker string  `json:"ticker"`
	Price  float64 `json:"price"`
}

type StockPricesBatch []StockPrice

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan StockPricesBatch)
var upgrader = websocket.Upgrader{}
var lock sync.Mutex

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port :%s...", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var prices StockPricesBatch
		err := ws.ReadJSON(&prices)
		if err != nil {
			log.Printf("readJSON error: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- prices
	}
}

func handleMessages() {
	for {
		prices := <-broadcast
		for client := range clients {
			err := client.WriteJSON(prices)
			if err != nil {
				log.Printf("WriteJSON error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func SimulateRealTimeStockPriceGenerator() {
	tickers := []string{"GOOG", "AAPL", "MSFT", "AMZN", "FB"}
	batchSize := 5
	for {
		time.Sleep(10 * time.Second)
		var priceUpdates StockPricesBatch
		for i := 0; i < batchSize; i++ {
			selectedTicker := tickers[rand.Intn(len(tickers))]
			priceUpdate := StockPrice{
				Ticker: selectedTicker,
				Price:  rand.Float64()*1000 + 100,
			}
			priceUpdates = append(priceUpdates, priceUpdate)
		}

		broadcast <- priceUpdates
		log.Printf("Generated prices: %v", priceUpdates)
	}
}