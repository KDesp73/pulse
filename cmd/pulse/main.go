package main

import (
	"fmt"
	"log"
	"net/http"
	client "pulse/internal/client"
	"pulse/internal/config"
	"pulse/internal/database"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)


func main() {
	database.InitDB()

	var err error
    config.GlobalConfig, err = config.LoadConfig("mqtt.yml")
    if err != nil {
        panic(fmt.Sprintf("Could not load config: %v", err))
    }

    // Set up MQTT client options
    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tls://%s:%d", config.GlobalConfig.MQTT.Server, config.GlobalConfig.MQTT.Port))
    opts.SetClientID("MQTT Client in Go")
    opts.SetUsername(config.GlobalConfig.MQTT.Username)
    opts.SetPassword(config.GlobalConfig.MQTT.Password)

    // Set the MQTT callbacks
    opts.SetDefaultPublishHandler(client.MessagePubHandler)
    opts.OnConnect = client.ConnectHandler
    opts.OnConnectionLost = client.ConnectLostHandler

    // Create the MQTT client
    mqttClient := mqtt.NewClient(opts)

    // Connect to the MQTT client
    if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

    // Subscribe to the MQTT topic
    client.GenerateSubscription(config.GlobalConfig.MQTT.Topic, 1)(mqttClient)

    // Set up HTTP handlers for the dashboard and SSE events
    http.HandleFunc("/", client.DashboardHandler)
    http.HandleFunc("/events", client.SSEHandler)

    // Start the web server
    port := config.GlobalConfig.Web.Port
    if port == 0 {
        port = 8080
    }
	log.Printf("Web server started at http://localhost:%d", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
