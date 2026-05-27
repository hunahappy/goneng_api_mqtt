package main

import "time"

type SensorMessage struct {
	Topic     string
	Payload   string
	ReceivedAt time.Time
}