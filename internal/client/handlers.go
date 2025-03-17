package client

import (
	"fmt"
	"log"
	"net/http"
	"pulse/internal/config"
	"pulse/internal/database"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Channel to broadcast MQTT messages
var Broadcast = make(chan string)

// MQTT message handler callback when message is received
var MessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	// Send the message payload to the Broadcast channel
	Broadcast <- string(msg.Payload())
	// Insert message to the database for later use
	database.InsertMessageToDB(string(msg.Payload()))
}

// Called upon successful connection to the MQTT broker
var ConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to MQTT Broker")
	// Subscribe to the topic when connected
	if token := client.Subscribe("plant/data", 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

// Called when the connection to the MQTT broker is lost
var ConnectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err)
}

// SSEHandler to handle Server-Sent Events
func SSEHandler(w http.ResponseWriter, r *http.Request) {
	// Set the appropriate headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Keep the connection open
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	for {
		// Wait until we have new data
		message := <-Broadcast

		// Directly send the raw JSON message as SSE event
		fmt.Fprintf(w, "data: %s\n\n", message)
		flusher.Flush()
	}
}

// Serve the dashboard page
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, config.GlobalConfig.Web.Page)
}
