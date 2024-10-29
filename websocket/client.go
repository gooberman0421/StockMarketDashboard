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
		log.Println(".env file not found")
	}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
}

func (manager *ClientManager) addClient(client *Client) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.clients[client] = true
}

func (manager *ClientManager) removeClient(client *Client) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	if _, ok := manager.clients[client]; ok {
		delete(manager.clients, client)
	}
}

func (manager *ClientManager) broadcast(message string) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	for client := range manager.clients {
		err := client.Conn.WriteJSON(broadcastMessage{Message: message})
		if err != nil {
			log.Printf("Error broadcasting message to client: %v", err)
			manager.removeClient(client)
		}
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	client := &Client{Conn: conn}
	clientManager.addClient(client)
	defer func() {
		clientManager.removeClient(client)
		conn.Close()
	}()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
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
	log.Println("Starting server on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}