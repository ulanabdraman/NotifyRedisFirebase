package functions

import (
	"Organize/geozone"
	"Organize/models"
	"Organize/params"
	"context"
	"encoding/json"
	"firebase.google.com/go"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sync"
)

func SubscribeToChannel(client *redis.Client, app *firebase.App, db *pgxpool.Pool) {
	pubsub := client.Subscribe(context.Background(), "realtime")
	defer func(pubsub *redis.PubSub) {
		err := pubsub.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(pubsub)

	ch := pubsub.Channel()

	pubsub2 := client.Subscribe(context.Background(), "updgeo")
	defer func(pubsub2 *redis.PubSub) {
		err := pubsub2.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(pubsub2)

	chgeo := pubsub2.Channel()
	var TripCh [4000]chan bool
	for i := 0; i < 4000; i++ {
		TripCh[i] = make(chan bool)
	}

	var wg sync.WaitGroup
	for {
		select {
		case msg := <-chgeo:
			var msginf int
			err := json.Unmarshal([]byte(msg.Payload), &msginf)
			if err != nil {
				fmt.Println("Error unmarshaling JSON2:", err)
				continue
			}
			geozone.InitGeo(msginf, db)
		case msg := <-ch:
			var msginf []models.MsgAll

			err := json.Unmarshal([]byte(msg.Payload), &msginf)
			if err != nil {
				fmt.Println("Error unmarshaling JSON1:", err)
				continue
			}
			if msginf[0].Id == 43 {
				log.Println(msginf)
				continue
			}

			for _, m := range msginf {
				params.CheckSpeed(m, app, db)
				params.CheckIgn(m, app, db)
			}

			params.CheckGeo(msginf[len(msginf)-1], app, db)

			if msginf[len(msginf)-1].Pos.S > 0 && params.TripActiv(db, msginf[len(msginf)-1].Id) {
				//log.Println("Отпарвлено ", msginf[len(msginf)-1].Id)
				go params.CheckTrip(msginf[len(msginf)-1], app, db, TripCh[msginf[len(msginf)-1].Id])
				TripCh[msginf[len(msginf)-1].Id] <- true
			}
			wg.Wait()
		}
	}

}
