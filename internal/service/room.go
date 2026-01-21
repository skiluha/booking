package service

import (
	"booking/internal/model"
	"booking/internal/repository"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type RoomService struct {
	conn *pgx.Conn
}

func NewRoomService(conn *pgx.Conn) *RoomService {
	return &RoomService{conn: conn}
}

func (s *RoomService) CreateRoom(
	ctx context.Context,
	name string,
	capacity int,
) (int, error) {
	if name == "" {
		return 0, errors.New("name is empty")
	}
	if capacity == 0 {
		return 0, errors.New("invalid capacity")
	}
	return repository.CreateRoom(ctx, s.conn, name, capacity)
}

func (s *RoomService) GetRooms(ctx context.Context) (rooms []model.Room, err error) {
	return repository.GetRooms(ctx, s.conn)
}
