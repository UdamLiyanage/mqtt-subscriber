package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"syscall"
)

var mqClient MQTT.Client

func init() {
	connect()
	c := cron.New()
	_, err := c.AddFunc(os.Getenv("CRON_TIME"), publishMean)
	if err != nil {
		panic(err)
	}
	c.Start()
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	if msg.Topic() != os.Getenv("MQTT_PUBLISH_TOPIC") {
		writeData(msg.Payload())
	}
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
	mqClient = MQTT.NewClient(opts)
	if token := mqClient.Connect(); token.Wait() && token.Error() != nil {
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

func publishMean() {
	mqClient.Publish(os.Getenv("MQTT_PUBLISH_TOPIC"), 0, false, getMean())
}

/* Device Definition
	ph, serial
    temp, humidity, light, serial
	CSV
*/
