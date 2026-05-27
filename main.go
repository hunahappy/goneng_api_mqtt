package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	cfg, _ := LoadConfigMap("goneng_mqtt.conf")

	db, err := NewDB(
		cfg["db_host"].(string),
		cfg["db_port"].(string),
		cfg["db_user"].(string),
		cfg["db_password"].(string),
		cfg["db_name"].(string),
		"disable",
	)
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}
	defer db.Close()

	client := NewMQTTClient(
		db,
		cfg["mqtt_broker"].(string),
		cfg["mqtt_username"].(string),
		cfg["mqtt_password"].(string),
	)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("초기 연결 실패: %v", token.Error())
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	client.Disconnect(250)
	fmt.Println("프로그램 종료")
}
