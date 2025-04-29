package main

import (
	"context"
	"database/sql"
	"flex-frog-bot/bot"
	"flex-frog-bot/db"
	"flex-frog-bot/db/repositories"
	repos "flex-frog-bot/db/repositories/interfaces"
	"flex-frog-bot/server"
	services2 "flex-frog-bot/services"
	"flex-frog-bot/services/interfaces"
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

	imgRepo, chatRepo, userRepo, searchRepo := initRepositories(conn)
	imgSvc, chatSvc, userSvc, searchSvc := initServices(imgRepo, chatRepo, userRepo, searchRepo)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		bot.RunBot(ctx, imgSvc, chatSvc, userSvc)
	}()
	go func() {
		defer wg.Done()
		server.RunServer(ctx, searchSvc)
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

func initRepositories(conn *sql.DB) (repos.ImageRepository, repos.ChatRepository, repos.UserRepository, repos.SearchRepository) {
	imgRepo := repositories.NewImageRepository(conn)
	chatRepo := repositories.NewChatRepository(conn)
	userRepo := repositories.NewUserRepository(conn)
	searchRepo := repositories.NewSearchRepository(conn)
	return imgRepo, chatRepo, userRepo, searchRepo
}

func initServices(imgRepo repos.ImageRepository, chatRepo repos.ChatRepository, userRepo repos.UserRepository, searchRepo repos.SearchRepository) (interfaces.ImageService, interfaces.ChatService, interfaces.UserService, interfaces.SearchService) {
	imgSvc := services2.NewImageService(imgRepo)
	chatSvc := services2.NewChatService(chatRepo)
	userSvc := services2.NewUserService(userRepo)
	searchSvc := services2.NewSearchService(searchRepo)
	return imgSvc, chatSvc, userSvc, searchSvc
}
