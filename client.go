package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	clientIdPrefix = "tracker"
	publisherType  = "pub"
	consumerType   = "sub"
	qosLevel       = 1
)

var (
	client  mqtt.Client
	mqtturi *url.URL
)

func init() {
	uri, err := url.Parse(os.Getenv("MQTT_URL"))
	if err != nil {
		log.Fatal(err)
	}

	mqtturi = uri
}

//connect set an instance of mqtt client with options
func connect(clientIdPrefix string, uri *url.URL, clientType string) {
	if client != nil && client.IsConnected() {
		return
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	clientId := fmt.Sprintf("%s_%s", clientIdPrefix, clientType)
	opts.SetClientID(clientId)
	opts.SetCleanSession(false) //for persistent session (fifo queue)

	c := mqtt.NewClient(opts)
	token := c.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}

	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	log.Printf("successfully connected to %s for client %s\n", fmt.Sprintf("tcp://%s", uri.Host), clientId)
	client = c
}

//publish connects to mqtt broker and publishes raw request to topic
func publish(topic string, req []byte) {
	connect(clientIdPrefix, mqtturi, publisherType)

	token := client.Publish(topic, qosLevel, false, req)
	if err := token.Error(); err != nil {
		log.Printf("error while publishing %s\n", err)
	} else {
		log.Printf("published to topic [%s]\n%s\n", topic, req)
	}
}

//consume connects to mqtt broker and subscribes to topic then call handler
func consume(topic string) {
	connect(clientIdPrefix, mqtturi, consumerType)

	if token := client.Subscribe(topic, qosLevel, messageHandler); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	disconnect()
	log.Println("stopping consumer")
}

//disconnect after 300 ms
func disconnect() {
	client.Disconnect(300)
	log.Printf("disconnected")
}
