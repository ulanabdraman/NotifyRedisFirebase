package params

import (
	"Organize/models"
	"Organize/notiffunc"
	"context"
	"firebase.google.com/go"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

var PTt [4000]time.Time
var PT [4000]bool

func CheckTrip(msg models.MsgAll, app *firebase.App, db *pgxpool.Pool, ch chan bool) {
	fmt.Println("Канал начат", msg.Id)
	if !PT[msg.Id] {
		PT[msg.Id] = true

	} else {
		fmt.Println("Канал безуспешно", msg.Id)
		return
	}
	fmt.Println("Все прошло успешно")
	T := convertUnixToTime(int64(msg.T))

	rowsCreators, err := db.Query(context.Background(), "SELECT user_id FROM main_unit_creator where unit_id = $1", msg.Id)

	if err != nil {
		log.Fatal(err)
	}
	defer rowsCreators.Close()

	log.Println("rows_creators ", rowsCreators)
	for rowsCreators.Next() {

		fmt.Println("Найдены пользватели")
		var user int
		rowsCreators.Scan(&user)
		var repText, repName string
		rowsFirebase, err := db.Query(context.Background(), "SELECT token FROM main_firebase_tokens where user_id=$1", user)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("rows_firebase ", rowsFirebase)
		defer rowsFirebase.Close()
		var deviceToken []string
		for rowsFirebase.Next() {
			//fmt.Println("Найдены токены")
			var token string
			rowsFirebase.Scan(&token)
			if sliceContains(deviceToken, token) {
				deviceToken = append(deviceToken, token)
				repName = fmt.Sprintf("Объект %s начал поездку", GetName(db, msg.Id))
				repText = fmt.Sprintf("Объект начал поездку в %s", T.Format("15:04:05"))
				err = notiffunc.SendFirebaseNotification(app, db, token, msg.Id, 0, repName, repText)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
	log.Println(rowsCreators)
	fmt.Println("Проверка")
	for {

		select {
		case <-ch:
			fmt.Println("Ням-ням")
			PTt[msg.Id] = time.Now()

		case <-time.After(5 * time.Minute):
			log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! stopped", msg.Id)
			PT[msg.Id] = false
			rowsCreators2, err := db.Query(context.Background(), "SELECT user_id FROM main_unit_creator where unit_id = $1", msg.Id)
			if err != nil {
				log.Fatal(err)
			}
			defer rowsCreators2.Close()

			for rowsCreators2.Next() {
				var user int
				rowsCreators.Scan(&user)
				var repText, repName string
				rowsFirebase, err := db.Query(context.Background(), "SELECT token FROM main_firebase_tokens where user_id=$1", user)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("rows_firebase ", rowsFirebase)
				defer rowsFirebase.Close()
				var deviceToken []string
				for rowsFirebase.Next() {
					var token string
					rowsFirebase.Scan(&token)
					if sliceContains(deviceToken, token) {
						deviceToken = append(deviceToken, token)
						repName = fmt.Sprintf("Объект %s закончил поездку", GetName(db, msg.Id))
						repText = fmt.Sprintf("Объект начал поездку в %s и закончил поездку в %s",
							T.Format("02.01.2006 15:04:05"),           // Формат даты и времени для начала поездки
							PTt[msg.Id].Format("02.01.2006 15:04:05"), // Формат даты и времени для завершения поездки
						)
						err = notiffunc.SendFirebaseNotification(app, db, token, msg.Id, 0, repName, repText)
						if err != nil {
							log.Println(err)
							continue
						}
					}
				}
			}
			return
		}

	}

}
