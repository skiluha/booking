package main

import (
	"booking/internal/db"
	"booking/internal/handler"
	"booking/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	conString := os.Getenv("CONN_STRING")
	if conString == "" {
		log.Fatal("Empty config")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Empty port")
	}

	rootctx := context.Background()
	conn := db.ConnectWithRetry(conString)
	if conn == nil {
		log.Println("db connection if nil")
	}
	defer conn.Close(rootctx)
	err := db.RunAutoMigrate(rootctx, conn)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	ctx, stop := signal.NotifyContext(rootctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, booking"))
	})
	roomService := service.NewRoomService(conn)
	roomHandler := handler.NewRoomHandler(roomService)
	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			roomHandler.GetRooms(w, r)
		case http.MethodPost:
			roomHandler.CreateRoom(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/rooms/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			roomHandler.GetRoombyId(w, r)
		case http.MethodDelete:
			roomHandler.DeleteRoom(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	go func() {
		fmt.Println("HTTP server stardet on port:", port)
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
