package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type Client struct {
	Conn *websocket.Conn
}

type ClientManager struct {
	clients map[*Client]bool
	lock    sync.Mutex
}

type broadcastMessage struct {
	Message string `json:"message"`
}

var (
	upgrader      = websocket.Upgrader{}
	clientManager = ClientManager{
		clients: make(map[*Client]bool),
	}
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found. Proceeding without it. Error: %v\n", err)
	}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
}

func (manager *ClientManager) addClient(client *Client) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.clients[client] = true
	log.Printf("Added new client. Total clients: %d\n", len(manager.clients))
}

func (manager *ClientManager) removeClient(client *Client) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	if _, ok := manager.clients[client]; ok {
		delete(manager.clients, client)
		log.Printf("Removed client. Total clients: %d\n", len(manager.clients))
	}
}

func (manager *ClientManager) broadcast(message string) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	for client := range manager.clients {
		err := client.Conn.WriteJSON(broadcastMessage{Message: message})
		if err != nil {
			log.Printf("Error broadcasting message to a client. Removing client. Error: %v\n", err)
			manager.removeClient(client)
			client.Conn.Close()
		}
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v\n", err)
		http.Error(w, "Could not upgrade to websocket", http.StatusInternalServerError)
		return
	}
	client := &Client{Conn: conn}
	clientManager.addClient(client)
	defer func() {
		clientManager.removeClient(client)
		conn.Close()
		log.Println("Closed connection for client.")
	}()
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			log.Printf("Error reading message from client. Error: %v\n", err)
			break
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/ws", handleConnections)
	go func() {
		for {
			clientManager.broadcast("Hello, World!")
			time.Sleep(10 * time.Second)
		}
	}()
	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}