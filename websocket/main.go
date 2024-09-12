package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
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

var (
	clients   = make(map[*websocket.Conn]bool) // Connected clients
	broadcast = make(chan StockPricesBatch)    // Broadcast channel
	upgrader  = websocket.Upgrader{}           // Configure the upgrader
	lock      sync.Mutex                        // Synchronize access
)

func init() {
	// Load .env file at the beginning
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()
	go SimulateRealTimeStockPriceGenerator() // Start the simulation

	// Determine port and start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port :%s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
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
		if err := ws.ReadJSON(&prices); err != nil {
			log.Printf("readJSON error: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- prices
	}
}

func handleMessages() {
	for prices := range broadcast {
		for client := range clients {
			if err := client.WriteJSON(prices); err != nil {
				log.Printf("WriteJSON error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func SimulateRealTimeStockPriceGenerator() {
	tickers := []string{"GOOG", "AAPL", "MSFT", "AMZN", "FB"}
	batchSize := len(tickers)
	for {
		time.Sleep(10 * time.Second) // Simulate delay
		var priceUpdates StockPricesBatch
		for _, ticker := range tickers {
			priceUpdate := StockPrice{
				Ticker: ticker,
				Price:  rand.Float64()*1000 + 100, // Generate a new price
			}
			priceUpdates = append(priceUpdates, priceUpdate)
		}

		broadcast <- priceUpdates
		log.Printf("Generated prices: %v", priceUpdates)
	}
}