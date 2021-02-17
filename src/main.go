package main

//importing necessary packages
import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

//initialize clients and broadcast channel
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

// Upgrader will allow us to turn a normal http connection into a WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Define our message object
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

//main entry point for this go application
func main() {
	// Create file server so users can access app.js and style.css
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// Handling requests for initiating a WebSocket
	http.HandleFunc("/ws", handleConnections)

	// Asynchronous call to listen to all messages
	go handleMessages()

	// Server will start on local and give any errors
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade http request to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Close connection when function is returned
	defer ws.Close()

	// Add new client to clients dictionary
	clients[ws] = true

  // Listen for any message and publish it to broadcast channel
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

//Grab messages from broadcast channel and relays message
func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
