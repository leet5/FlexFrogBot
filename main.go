package main

import (
	"context"
	"database/sql"
	"flex-frog-bot/bot"
	"flex-frog-bot/db"
	"flex-frog-bot/db/repository"
	tgbotapi "flex-frog-bot/tg-bot-api"
	_ "fmt"
	"log"
)

func main() {
	updates := tgbotapi.GetUpdatesChan()

	conn, err := db.NewPostgresConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			log.Printf("[database] Error closing connection: %v", err)
		}
	}(conn)

	bot.ImageRepo = repository.NewImageRepository(conn)
	bot.ChatRepo = repository.NewChatRepository(conn)
	bot.UserRepo = repository.NewUserRepository(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.ProcessUpdates(ctx, updates)
}
