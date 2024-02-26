package notiffunc

import (
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strconv"
)

func SendFirebaseNotification(app *firebase.App, db *pgxpool.Pool, deviceToken string, Id int, notif_id int, notif_name string, notif_text string) error {
	client, err := app.Messaging(context.Background())
	if err != nil {
		return err
	}
	// Определите сообщение для отправки
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: notif_name,
			Body:  notif_text,
		},
		Data: map[string]string{
			"notify_id": strconv.Itoa(notif_id),
			"unit_id":   strconv.Itoa(Id),
			// Добавьте любые другие данные, которые вы хотите включить в уведомление
		},
		Token: deviceToken,
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
	}

	// Отправка сообщения
	_, err = client.Send(context.Background(), message)
	if err != nil {
		return err
	} else {
		log.Println("Успешное отправление")
	}

	return nil
}
