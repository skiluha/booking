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
	defer rows.Close()
	for rows.Next() {
		var room model.Room
		err = rows.Scan(
			&room.ID,
			&room.Name,
			&room.Capacity,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rooms, err
}

func GetRoombyId(
	ctx context.Context,
	conn *pgx.Conn,
	id int,
) (model.Room, error) {

	sqlQuery := `
SELECT id, name, capacity
FROM rooms
WHERE id = $1
`
	var room model.Room

	err := conn.QueryRow(ctx, sqlQuery, id).Scan(&room.ID, &room.Name, &room.Capacity)
	if err != nil {
		return model.Room{}, err
	}
	return room, nil
}

func DeleteRooms(
	ctx context.Context,
	conn *pgx.Conn,
	id int,
) error {
	sqlQuery := `
DELETE FROM rooms
WHERE id = $1`
	_, err := conn.Exec(ctx, sqlQuery, id)
	if err != nil {
		return err
	}
	return nil
}
