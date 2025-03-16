package main

import (
	"fmt"
	"log"
	broker "mqtt-client/internal/broker"
	"mqtt-client/internal/config"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)


func main() {
	conf, err := config.LoadConfig("mqtt.yml")
	if err != nil {
		panic("Could not load config")
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", conf.MQTT.Server, conf.MQTT.Port))
	opts.SetClientID("MQTT Client in Go")
	opts.SetUsername(conf.MQTT.Username)
	opts.SetPassword(conf.MQTT.Password)

	opts.SetDefaultPublishHandler(broker.MessagePubHandler)
	opts.OnConnect = broker.ConnectHandler
	opts.OnConnectionLost = broker.ConnectLostHandler

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	broker.GenerateSubscription(conf.MQTT.Topic, 1)(client)

	http.HandleFunc("/", broker.DashboardHandler)

	port := conf.Web.Port
	if port == 0 {
		port = 8080
	}
	log.Printf("Web server started on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

