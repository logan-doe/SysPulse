package services

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"syspulse/internal/models"

	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	clients   map[*websocket.Conn]bool  // all connected clients
	broadcast chan models.SystemMetrics // channel for sending metrics
	mu        sync.Mutex                // anti race
	upgrader  websocket.Upgrader        // to upgrade http to websocket
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan models.SystemMetrics),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// В разработке разрешаем все origin
				// В продакшене нужно ограничить домены
				return true
			},
		},
	}
}

func (ws *WebSocketService) Start() {
	go ws.handleBroadcast() // goroutine for send messages to all clients
}

func (ws *WebSocketService) HandleConnection(w http.ResponseWriter, r *http.Request) {
	//upgrade http to web socket
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ Failed to upgrade http to web socket: %v", err)
		return
	}

	defer conn.Close()

	// New client
	ws.registerClients(conn)
	defer ws.unregisterClient(conn)

	log.Printf("🔌 New Web Socket client is connected: %s\n", r.RemoteAddr)

	// infinite cycle for reading messages from client
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("🔌 client disconnected: %s\n", r.RemoteAddr)
			break
		}
	}
}

func (ws *WebSocketService) registerClients(conn *websocket.Conn) { // add client to connections list
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.clients[conn] = true
}

func (ws *WebSocketService) unregisterClient(conn *websocket.Conn) { // remove client from connections list
	ws.mu.Lock()
	defer ws.mu.Unlock()
	delete(ws.clients, conn)
}

func (ws *WebSocketService) BroadcastMessage(metrics models.SystemMetrics) { // send metrics to all clients
	ws.broadcast <- metrics
}

func (ws *WebSocketService) handleBroadcast() { // processing message sending
	for metrics := range ws.broadcast {
		message, err := json.Marshal(metrics)
		if err != nil {
			log.Printf("❌ Failed to marshal metric: %v", err)
			continue
		}
		ws.sendToAllCLients(message)
	}
}

func (ws *WebSocketService) sendToAllCLients(message []byte) { // sending message to each and every connected client
	ws.mu.Lock()
	defer ws.mu.Unlock()

	for client := range ws.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf(" Failed to send Web Socket message: %v", err)
			client.Close()
			delete(ws.clients, client)
		}
	}
}

func (ws *WebSocketService) GetConnectedClientCount() int { // returning number of al connected clients
	ws.mu.Lock()
	defer ws.mu.Unlock()
	return len(ws.clients)
}
