package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func CreateBooking(
	ctx context.Context,
	conn *pgx.Conn,
	roomID int,
	startDate, endDate time.Time,
) (int, error) {
	sqlQuery := `
INSERT INTO bookings (room_id, start_date, end_date)
VALUES ($1,$2,$3),
RETURNING ID;
`
	var id int
	err := conn.QueryRow(ctx, sqlQuery, roomID, startDate, endDate).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
