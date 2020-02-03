package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	connect()
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	writeData(msg.Payload())
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	opts := setClientOptions()
	topic := "#"

	opts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server:\n")
	}
	<-c
}

func setClientOptions() *MQTT.ClientOptions {
	opts := MQTT.NewClientOptions().AddBroker(os.Getenv("MQTT_BROKER"))
	opts.SetDefaultPublishHandler(f)
	return opts
}

/* Device Definition
	ph, serial
    temp, humidity, light, serial
	CSV
*/
