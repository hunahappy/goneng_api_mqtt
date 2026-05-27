package main

import (
	"database/sql"
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

func InsertSensorMessage(db *sql.DB, msg SensorMessage) error {
	query := `
		INSERT INTO 센서 (토픽, 내용)
		VALUES ($1, $2)
	`
	_, err := db.Exec(query, msg.Topic, msg.Payload)
	return err
}