package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type StockPrice struct {
	Ticker string  `json:"ticker"`
	Price  float64 `json:"price"`
}

var clients = make(map[*websocket.Conn]bool)

var broadcast = make(chan StockPrice)

var upgrader = websocket.Upgrader{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
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
		return // Return instead of fatal to not kill server
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var price StockPrice
		err := ws.ReadJSON(&price)
		if err != nil {
			log.Printf("readJSON error: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- price
	}
}

func handleMessages() {
	for {
		price := <-broadcast
		for client := range clients {
			err := client.WriteJSON(price)
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
	for {
		time.Sleep(time.Duration(rand.Intn(10)+1) * time.Second)

		selectedTicker := tickers[rand.Intn(len(tickers))]
		priceUpdate := StockPrice{
			Ticker: selectedTicker,
			Price:  rand.Float64()*1000 + 100, // Simulate price between 100 and 1100
		}

		broadcast <- priceUpdate
		log.Printf("Generated price: %v", priceUpdate)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
	log.SetOutput(os.Stdout) // Make sure logging goes to stdout
	go SimulateRealTimeStockPriceGenerator()
}