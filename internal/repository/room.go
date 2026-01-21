package repository

import (
	"booking/internal/model"
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateRoom(
	ctx context.Context,
	conn *pgx.Conn,
	name string,
	capacity int,
) (int, error) {

	sqlQuery := `
	INSERT INTO rooms (name, capacity)
	VALUES ($1, $2)
	RETURNING id
`
	var id int
	err := conn.QueryRow(ctx, sqlQuery, name, capacity).Scan(&id)
	return id, err
}

func GetRooms(
	ctx context.Context,
	conn *pgx.Conn,
) ([]model.Room, error) {
	sqlQuery := `
SELECT id, name, capacity
FROM rooms
ORDER BY id ASC
`
	rows, err := conn.Query(ctx, sqlQuery)
	rooms := make([]model.Room, 0)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var room model.Room
		err = rows.Scan(
			&room.ID,
			&room.Name,
			&room.Capacity,
		)
		if err != nil {
			panic(err)
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rooms, err

}
