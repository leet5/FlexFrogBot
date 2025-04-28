package main

import (
	"context"
	"database/sql"
	"flex-frog-bot/bot"
	"flex-frog-bot/db"
	"flex-frog-bot/db/repository"
	"flex-frog-bot/server"
	_ "fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
		bot.RunBot(ctx, imgRepo, chatRepo, userRepo)
	}()
	go func() {
		defer wg.Done()
		server.RunServer(ctx, imgRepo)
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
