package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)

func NewDB(host, port, user, password, dbname string, sslmode string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	} else {
		fmt.Println("postgres 연결 성공")
	}

	return db, nil
}

func InsertSensorMessage(db *sql.DB, data map[string]interface{}) error {
	query := `
		INSERT INTO 로그 (장치, 구분, 내용, 토픽)
		VALUES ($1, $2, $3, $4)
	`

	// 안전하게 문자열로 변환
	device := fmt.Sprintf("%v", data["장치"])
	category := fmt.Sprintf("%v", data["구분"])
	topic := fmt.Sprintf("%v", data["토픽"])

	contentBytes, err := json.Marshal(data["내용"])
	if err != nil {
		return err
	}
	content := string(contentBytes)

	_, err = db.Exec(query, device, category, content, topic)
	fmt.Printf("DB에 저장된 데이터: 장치=%s, 구분=%s, 토픽=%s, 내용=%s\n", device, category, topic, content)
	return err
}
