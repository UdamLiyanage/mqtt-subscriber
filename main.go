package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/signal"
	"syscall"
)

//var knt int
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	opts := setClientOptions()
	topic := "test/button"

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
	opts := MQTT.NewClientOptions().AddBroker("tcp://54.67.77.75:1883")
	opts.SetDefaultPublishHandler(f)
	return opts
}
