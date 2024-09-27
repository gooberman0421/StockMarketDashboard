package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust this to a proper origin checking for production
	},
}

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file, %s\n", err)
	}

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set.")
	}

	log.Println("http server started on :" + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("ListenAndServe failed: %v\n", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v\n", err)
		return // Don't call log.Fatal here, otherwise server will stop on error
	}
	defer func() {
		err := ws.Close()
		if err != nil {
			log.Printf("error closing connection: %v\n", err)
		}
	}()

	clients[ws] = true

	for {
		var msg Message
		if err := ws.ReadJSON(&msg); err != nil {
			log.Printf("error reading json: %v\n", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for msg := range broadcast {
		for client := range clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Printf("error writing json: %v\n", err)
				_ = client.Close()
				delete(clients, client)
			}
		}
	}
}