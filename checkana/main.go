package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Cell)              // broadcast channel for the numbers slice
var cells = []Cell{}                         // initial slice of numbers

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

	cells = initialize10000Cells()

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
		cell := Cell{}
		err = json.Unmarshal(message, &cell)
		fmt.Println("Received", cell)
		if err != nil {
			log.Printf("error: invalid number received %v", err)
		} else {
			flipCellWithId(cell.CellID)
			broadcast <- cell
		}
	}
}

func handleMessages() {
	for {
		updatedNumbers := <-broadcast
		numbersJSON, err := json.Marshal(updatedNumbers)
		fmt.Println("Broadcasting", updatedNumbers)
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
	numbersJSON, err := json.Marshal(cells)
	if err != nil {
		log.Println("Error encoding slice of numbers:", err)
		return
	}
	ws.WriteMessage(websocket.TextMessage, numbersJSON)
}

func initialize10000Cells() []Cell {
	for i := 0; i < 10000; i++ {
		cells = append(cells, Cell{CellID: i, State: "🍌"})
	}
	return cells
}

func flipCellWithId(id int) {
	for i, cell := range cells {
		if cell.CellID == id {
			if cell.State == "🍌" {
				cells[i].State = "🍎"
			} else {
				cells[i].State = "🍌"
			}
		}
	}
}
