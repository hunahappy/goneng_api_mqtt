package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewMQTTClient(db *sql.DB, mqtt_broker, user, password string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqtt_broker)

	clientID := "device-" + uuid.New().String()
	opts.SetClientID(clientID)
	opts.SetUsername(user)
	opts.SetPassword(password)
	opts.SetCleanSession(false)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(10 * time.Second)

	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("MQTT 브로커 연결 성공")

		if token := c.Subscribe("goneng/farm1/#", 1, func(client mqtt.Client, msg mqtt.Message) {
			// handleMessage(db, client, msg)
			handleMessage(db, msg)
		}); token.Wait() && token.Error() != nil {
			log.Printf("구독 실패: %v", token.Error())
		} else {
			fmt.Println("구독 완료: goneng/farm1/#")
		}
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Printf("MQTT 연결 끊김: %v\n", err)
		fmt.Println("자동 재연결 대기 중...")
	}

	return mqtt.NewClient(opts)
}

// func handleMessage(db *sql.DB, client mqtt.Client, msg mqtt.Message) {
func handleMessage(db *sql.DB, msg mqtt.Message) {
	fmt.Printf("메시지 수신 - 주제: %s, 내용: %s\n", msg.Topic(), msg.Payload())

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(msg.Payload()), &data); err != nil {
		panic(err)
	}
	data["토픽"] = msg.Topic()

	if err := InsertSensorMessage(db, data); err != nil {
		log.Printf("DB 저장 실패: %v", err)
		return
	}
}
