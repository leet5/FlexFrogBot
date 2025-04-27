package main

import (
	"context"
	"database/sql"
	"errors"
	"flex-frog-bot/bot"
	"flex-frog-bot/db"
	"flex-frog-bot/db/repository"
	"flex-frog-bot/server"
	tgbotapi "flex-frog-bot/tg-bot-api"
	_ "fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	conn, err := db.NewPostgresConnection()
	if err != nil {
		log.Fatalf("[main] ❌ Failed to connect to database: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("[main] ⚠️ Error closing database connection: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handleShutdown(cancel)

	imgRepo, chatRepo, userRepo := initRepositories(conn)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		runBot(ctx, imgRepo, chatRepo, userRepo)
	}()
	go func() {
		defer wg.Done()
		runServer(ctx, imgRepo)
	}()

	wg.Wait()
	log.Println("[main] ✅ Application stopped gracefully")
}

func handleShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Println("[main] ⚠️ Received shutdown signal, stopping...")
	cancel()
}

func initRepositories(conn *sql.DB) (*repository.ImageRepository, *repository.ChatRepository, *repository.UserRepository) {
	imgRepo := repository.NewImageRepository(conn)
	chatRepo := repository.NewChatRepository(conn)
	userRepo := repository.NewUserRepository(conn)
	return imgRepo, chatRepo, userRepo
}

func runBot(ctx context.Context, imgRepo *repository.ImageRepository, chatRepo *repository.ChatRepository, userRepo *repository.UserRepository) {
	bot.ImgRepo = imgRepo
	bot.ChatRepo = chatRepo
	bot.UserRepo = userRepo

	updates := tgbotapi.GetUpdatesChan()
	bot.ProcessUpdates(ctx, updates)
}

func runServer(ctx context.Context, imgRepo *repository.ImageRepository) {
	server.ImgRepo = imgRepo

	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	http.HandleFunc("/api/v1/search", server.WithCORS(server.HandleSearch))

	go func() {
		log.Println("[server] ✅ Starting server on :8080 ...")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("[server] ⚠️ Error starting server: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("[server] ⚠️ Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("[server] ❌ Error during server shutdown: %v", err)
	}
}
