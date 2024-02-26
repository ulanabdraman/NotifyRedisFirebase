package params

import (
	"Organize/models"
	"Organize/notiffunc"
	"context"
	"firebase.google.com/go"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var PI []int

func CheckIgn(msg models.MsgAll, app *firebase.App, db *pgxpool.Pool) {
	//log.Println("ignition check ", msg.Id)
	T := convertUnixToTime(int64(msg.T))
	rows_notif_unit, err := db.Query(context.Background(), "SELECT * FROM main_notification_units where unit_id = $1", msg.Id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows_notif_unit.Close()

	for rows_notif_unit.Next() {
		var id, notif_id, unit_id int
		rows_notif_unit.Scan(&id, &notif_id, &unit_id)
		rows_notif, err := db.Query(context.Background(), "SELECT name,activated,id,param,creator_id from main_notification where id=$1", notif_id)
		if err != nil {
			log.Fatal(err)
		}

		var rep models.NotifParam
		var Max []byte
		for rows_notif.Next() {
			err = rows_notif.Scan(&rep.Name, &rep.Activated, &rep.Param_id, &Max, &rep.Creator_id)
			if err != nil {
				log.Fatal(err)
			}

			if rep.Activated && rep.Param_id == 3 {

				log.Println("зажигание", T, msg.Pos.S)

				rows_firebase, err := db.Query(context.Background(), "SELECT token FROM main_firebase_tokens where user_id=$1", rep.Creator_id)
				if err != nil {
					log.Fatal(err)
				}
				defer rows_firebase.Close()

				var deviceToken []string
				for rows_firebase.Next() {
					var token string
					rows_firebase.Scan(&token)
					if sliceContains(deviceToken, token) {
						deviceToken = append(deviceToken, token)
						var repText string
						if msg.Ign > 0 {
							repText = fmt.Sprintf("Включение зажигания у объекта %d", msg.Id)
						} else {
							repText = fmt.Sprintf("Отключение зажигания у объекта %d", msg.Id)
						}
						err = notiffunc.SendFirebaseNotification(app, db, token, msg.Id, notif_id, repText, rep.Name)
						if err != nil {
							log.Println(err)
							continue
						}
					}
				}
			}
		}
	}
}
