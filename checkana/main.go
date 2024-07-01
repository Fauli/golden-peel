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

type Emojis struct {
	Emojis []string `json:"emojis"`
}

func main() {

	// WE NEED A LOT OF BANANAS!!
	cells = initializeCells(10000)

	http.HandleFunc("/ws", handleConnections)
	http.Handle("/", http.FileServer(http.Dir("./grid")))
	http.Handle("/initial", handleInitial())

	go handleMessages()

	log.Println("HTTP server started on :9000")
	err := http.ListenAndServe("localhost:9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// handleConnections upgrades the HTTP connection to a WebSocket connection
func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Keep track of them clients! :)
	clients[ws] = true

	// Send the initial slice of cells to the client upon connection
	sendCells(ws)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Convert message to cell, set it, and send it to all clients
		cell := Cell{}
		err = json.Unmarshal(message, &cell)
		fmt.Println("Received", cell)
		if err != nil {
			log.Printf("error: invalid number received %v", err)
		} else {
			if len(cell.State) == 1 {
				setCellWithId(cell.CellID, cell.State)
				broadcast <- cell
			} else {
				log.Printf("error: invalid state received [%v], will ignore", err)
			}
		}
	}
}

// handleMessages broadcasts the updated numbers to all connected clients
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

// Helper function to send the updated cell
func sendCells(ws *websocket.Conn) {
	numbersJSON, err := json.Marshal(cells)
	if err != nil {
		log.Println("Error encoding slice of numbers:", err)
		return
	}
	ws.WriteMessage(websocket.TextMessage, numbersJSON)
}

// initializeCells creates a slice of cells of BANANAS with the given count
func initializeCells(count int) []Cell {
	for i := 0; i < count; i++ {
		cells = append(cells, Cell{CellID: i, State: "🍌"})
	}
	return cells
}

// TODO: [franz] This is like... totally secure and stuff :P
func setCellWithId(id int, state string) {
	for i, cell := range cells {
		if cell.CellID == id {
			cells[i].State = state
		}
	}
}

// handleInitial returns the initial state of the grid
func handleInitial() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Emojis{Emojis: cellsToEmojis(cells)})
	})
}

// cellsToEmojis converts a slice of cells to a slice of emojis
// TODO: [franz] This is a bit of a hack, but it works for now
func cellsToEmojis(cells []Cell) []string {
	emojis := []string{}
	for _, cell := range cells {
		emojis = append(emojis, cell.State)
	}
	return emojis
}
