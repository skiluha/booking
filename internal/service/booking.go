package service

import (
	"booking/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type BookingService struct {
	conn *pgx.Conn
}

func NewBookingService(conn *pgx.Conn) *BookingService {
	return &BookingService{conn: conn}
}

func (s *BookingService) CreateBooking(
	ctx context.Context,
	roomID int,
	startDate, endDate time.Time) (int, error) {
	if endDate.Before(startDate) {
		return 0, errors.New("invalid date")
	}

	sql := `
	SELECT COUNT(*) 
FROM booking
	WHERE room_id = $1
		AND start_date<$2
		AND end_date>$3
`
	var count int
	err := s.conn.QueryRow(ctx, sql, roomID, endDate, startDate).Scan(&count)
	if err != nil {
		return 0, err
	}

	if count > 0 {
		return 0, errors.New("room is already booked for this period")
	}
	return repository.CreateBooking(ctx, s.conn, roomID, startDate, endDate)
}
