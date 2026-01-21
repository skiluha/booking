package main

import (
	"booking/internal/handler"
	"booking/internal/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	conString := os.Getenv("CON_STRING")
	if conString == "" {
		log.Fatal("Empty config")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Empty port")
	}

	rootctx := context.Background()
	conn, err := pgx.Connect(rootctx, conString)
	if err != nil {
		panic(err)
	}
	defer conn.Close(rootctx)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, booking"))
	})
	ctx, stop := signal.NotifyContext(rootctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	roomService := service.NewRoomService(conn)
	roomHandler := handler.NewRoomHandler(roomService)
	mux.HandleFunc("/rooms", roomHandler.CreateRoom)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("server shutdown error:", err)
	}
}
