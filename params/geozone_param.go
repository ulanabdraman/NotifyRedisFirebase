package params

import (
	"Organize/geozone"
	"Organize/models"
	"Organize/notiffunc"
	"context"
	"firebase.google.com/go"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func CheckGeo(msg models.MsgAll, app *firebase.App, db *pgxpool.Pool) {
	log.Println("geozone check ", msg.Id)
	T := convertUnixToTime(int64(msg.T))
	var notify_geozone bool
	err := db.QueryRow(context.Background(), "SELECT notify_geozones FROM main_unit where id = $1", msg.Id).Scan(&notify_geozone)
	if err != nil {
		log.Fatal(err)
	}
	rows_creators, err := db.Query(context.Background(), "SELECT user_id FROM main_unit_creator where unit_id = $1", msg.Id)
	if err != nil {
		log.Fatal(err)
	}
	if notify_geozone {
		for rows_creators.Next() {
			var user int
			rows_creators.Scan(&user)
			change, inout, name := geozone.FindPolygon(user, models.Pos{msg.Pos.X, msg.Pos.Y, 0}, msg.Id)
			change2, inout2, name2 := geozone.FindCircle(user, models.Pos{msg.Pos.X, msg.Pos.Y, 0}, msg.Id)
			change = change || change2
			inout = append(inout, inout2...)
			name = append(name, name2...)
			var repText, repName string
			if change {
				rows_firebase, err := db.Query(context.Background(), "SELECT token FROM main_firebase_tokens where user_id=$1", user)
				if err != nil {
					log.Fatal(err)
				}
				defer rows_firebase.Close()
				var deviceToken []string
				for rows_firebase.Next() {
					var token string
					rows_firebase.Scan(&token)
					if sliceContains(deviceToken, token) {
						for i, _ := range inout {
							if inout[i] {
								repText = fmt.Sprintf("Объект %d заехал в геозону %s в %t", msg.Id, name[i], T)
								repName = fmt.Sprintf("Объект %d заехал в геозону %s", msg.Id, name[i])
							} else {
								repText = fmt.Sprintf("Объект %d покинул геозону %s в %t", msg.Id, name[i], T)
								repName = fmt.Sprintf("Объект %d покинул геозону %s", msg.Id, name[i])
							}
							err = notiffunc.SendFirebaseNotification(app, db, token, msg.Id, 0, repText, repName)
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

}
