package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)       // connected clients
var broadcast = make(chan []int)                   // broadcast channel for the numbers slice
var numbers = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} // initial slice of numbers

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

type Cell struct {
	CellID int    `json:"cell"`
	State  string `json:"state"`
}

func main() {

	http.HandleFunc("/ws", handleConnections)
	http.Handle("/", http.FileServer(http.Dir("./grid")))

	go handleMessages()

	log.Println("HTTP server started on :8000")
	err := http.ListenAndServe("localhost:8000", nil)
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

	// Send the initial slice of numbers to the client upon connection
	sendNumbers(ws)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Convert message to an integer and add to slice
		newNumber, err := strconv.Atoi(string(message))
		if err != nil {
			log.Printf("error: invalid number received %v", err)
		} else {
			numbers = append(numbers, newNumber)
			broadcast <- numbers
		}
	}
}

func handleMessages() {
	for {
		updatedNumbers := <-broadcast
		numbersJSON, err := json.Marshal(updatedNumbers)
		if err != nil {
			log.Printf("Error marshaling numbers: %v", err)
			continue
		}
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, numbersJSON)
			if err != nil {
				log.Printf("websocket error: %s", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// Helper function to send the current slice of numbers to a specific WebSocket client
func sendNumbers(ws *websocket.Conn) {
	numbersJSON, err := json.Marshal(numbers)
	if err != nil {
		log.Println("Error encoding slice of numbers:", err)
		return
	}
	ws.WriteMessage(websocket.TextMessage, numbersJSON)
}
