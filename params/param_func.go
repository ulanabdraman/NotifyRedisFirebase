package params

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strings"
	"time"
)

func convertUnixToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func sliceContains(slice []string, element string) bool {
	for _, value := range slice {
		if value == element {
			return false
		}
	}
	return true
}
func GetName(db *pgxpool.Pool, Id int) string {
	var name string
	err := db.QueryRow(context.Background(), "Select name from main_unit where id=$1", Id).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	return strings.ToUpper(name)
}
func TripActiv(db *pgxpool.Pool, Id int) bool {
	var activ bool
	err := db.QueryRow(context.Background(), "Select notify_trips from main_unit where id=$1", Id).Scan(&activ)
	if err != nil {
		log.Fatal(err)
	}
	return activ
}
