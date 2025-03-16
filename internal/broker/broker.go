package broker

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func GenerateSubscription(topic string, qos byte) func(client mqtt.Client) {
	return func(client mqtt.Client) {
		token := client.Subscribe(topic, qos, nil)
		token.Wait()
		if token.Error() != nil {
			fmt.Printf("Failed to subscribe to topic %s: %v\n", topic, token.Error())
			panic(token.Error())
		}
		fmt.Printf("Subscribed to topic: %s\n", topic)
	}
}

func GeneratePublication(topic string, msg string, qos byte, retained bool) func(client mqtt.Client) {
	return func(client mqtt.Client) {
		token := client.Publish(topic, qos, retained, msg)
		token.Wait()
		if token.Error() != nil {
			fmt.Printf("Failed to publish to topic %s: %v\n", topic, token.Error())
			panic(token.Error())
		}
		fmt.Printf("Published message to topic: %s\n", topic)
	}
}
