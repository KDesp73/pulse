package broker

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan string) // Could also use a struct for JSON data

var upgrader = websocket.Upgrader{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
    upgrader.CheckOrigin = func(r *http.Request) bool { return true } // Allow all origins (for development)

    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer ws.Close()

    clients[ws] = true

    for {
        _, _, err := ws.ReadMessage()
        if err != nil {
            delete(clients, ws)
            break
        }
    }
}

func handleMessages() {
    for {
        msg := <-Broadcast
        for client := range clients {
            err := client.WriteMessage(websocket.TextMessage, []byte(msg))
            if err != nil {
                client.Close()
                delete(clients, client)
            }
        }
    }
}
