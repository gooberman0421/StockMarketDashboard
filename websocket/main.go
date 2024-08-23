package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

	log.Println("Starting server on port :8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var price StockPrice
		err := ws.ReadJSON(&price)
		if err != nil {
			log.Printf("error: %v", err)
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
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func SimulateRealTimeStockPriceGenerator() {
	for {
		time.Sleep(10 * time.Second)

		priceUpdate := StockPrice{
			Ticker: "GOOG",
			Price:  2342.05,
		}

		broadcast <- priceUpdate
	}
}

func init() {
	go SimulateRealTimeStockPriceGenerator()
}