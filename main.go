package main

import (
	"Organize/functions"
	"Organize/geozone"
	_ "github.com/lib/pq" // PostgreSQL driver

	"Organize/databases"
	"log"
	"time"
)

func main() {
	// Подключение к Redis
	client := databases.ConnectToRedis()

	db := databases.ConnectToPostgres()

	geozone.StartSave(db)

	defer db.Close()
	// Подключение к Firebase
	app, err := databases.ConnectToFirebase()
	if err != nil {
		log.Fatal(err)
	}

	// Подписка на канал в горутине
	functions.SubscribeToChannel(client, app, db)

	// Ожидание завершения работы подписчика
	time.Sleep(time.Millisecond * 500)
}
