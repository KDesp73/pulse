package broker

import (
	"fmt"
	"html/template"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// this callback triggers when a message is received, it then prints the message (in the payload) and topic
var MessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	Broadcast <- string(msg.Payload())
}

// upon connection to the client, this is called
var ConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

// this is called when the connection to the client is lost, it prints "Connection lost" and the corresponding error
var ConnectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	type DashboardData struct {
		Temperature string
		Humidity    string
		Soil        string
		Light       string
	}

	data := DashboardData{
		Temperature: "-- °C",
		Humidity:    "-- %",
		Soil:        "-- %",
		Light:       "-- %",
	}

	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Plant Monitor Dashboard</title>
		<style>
			body { font-family: Arial, sans-serif; background-color: #f9f9f9; padding: 2rem; }
			h1 { color: #333; }
			.card { background: #fff; padding: 1rem; margin: 1rem 0; border-radius: 8px; box-shadow: 0 0 10px rgba(0,0,0,0.1); }
		</style>
	</head>
	<body>
		<h1>🌱 Plant Monitor Dashboard</h1>
		<div class="card">
			<p>Temperature: {{.Temperature}}</p>
			<p>Humidity: {{.Humidity}}</p>
			<p>Soil Moisture: {{.Soil}}</p>
			<p>Light: {{.Light}}</p>
		</div>
	</body>
	</html>`

	t, err := template.New("dashboard").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

