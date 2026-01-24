package db

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func RunAutoMigrate(ctx context.Context, conn *pgx.Conn) error {
	if conn == nil {
		return errors.New("nil connection")
	}
	sql := `
CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    capacity INT NOT NULL
);

CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    room_id INT NOT NULL,
	start_date DATE NOT NULL,
	end_date DATE NOT NULL
);
`
	_, err := conn.Exec(ctx, sql)
	return err
}

func ConnectWithRetry(conn_string string) *pgx.Conn {

	var conn *pgx.Conn
	var err error

	for i := 0; i <= 10; i++ {
		conn, err = pgx.Connect(context.Background(), conn_string)
		if err == nil {
			log.Println("database connected")
			return conn
		}

		log.Printf("database not ready (попытка %d),%v", i, err)

		time.Sleep(2 * time.Second)

		return nil
	}
	log.Fatal("failed to connect to database")
	return nil
}
