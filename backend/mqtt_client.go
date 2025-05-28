package main

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func StartMQTT() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("go_dashboard")
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	client.Subscribe("esp32/pzem", 0, func(c mqtt.Client, m mqtt.Message) {
		var data PZEMPayload
		if err := json.Unmarshal(m.Payload(), &data); err == nil {
			InsertPZEMData(data)
		}
	})

	client.Subscribe("esp32/xymd", 0, func(c mqtt.Client, m mqtt.Message) {
		var data XYMDPayload
		if err := json.Unmarshal(m.Payload(), &data); err == nil {
			InsertXYMDData(data)
		}
	})
}
