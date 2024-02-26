// postgres_connection.go
package databases

import (
	"context"
	//"database/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
)

func ConnectToPostgres() *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), "postgres://erema:loc25387@localhost:5432/glot")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	//err = db.Ping()
	//if err != nil {
	//	log.Fatal(err)
	//}

	log.Println("Успешное подключение к postgres")

	return db
}
