package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	client "pulse/internal/client"
	"pulse/internal/config"
	"pulse/internal/database"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const configTemplate = `mqtt:
  server: "broker.mqtt"
  port: 8883
  username: "your-username"
  password: "your-password"
  topic: "your/topic"
web:
  port: 8080
  page: your-dashboard.html
`

func generateConfig() {
	file, err := os.Create("pulse.yml")
	if err != nil {
		fmt.Println("Error creating config file:", err)
		return
	}
	defer file.Close()

	tmpl, err := template.New("config").Parse(configTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = tmpl.Execute(file, nil)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Println("Config file generated successfully: config.yml")
}

func main() {
	generateConfigFlag := flag.Bool("generate-config", false, "Generate a template configuration file")

	flag.Parse()

	if *generateConfigFlag {
		generateConfig()
		return
	}

	database.InitDB()

	var err error
    config.GlobalConfig, err = config.LoadConfig("pulse.yml")
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
	http.HandleFunc("/api/avg-temp", database.GetAvgTemperature)
	http.HandleFunc("/api/min-max-temp", database.GetMinMaxTemperature)
	http.HandleFunc("/api/avg-moist", database.GetAvgSoilMoisture)
	http.HandleFunc("/api/latest", database.GetLatestReading)

    // Start the web server
    port := config.GlobalConfig.Web.Port
    if port == 0 {
        port = 8080
    }
	log.Printf("Web server started at http://localhost:%d", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
