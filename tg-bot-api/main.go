package tg_bot_api

import (
	"log"
	"os"
)

var (
	token string
)

func init() {
	token = os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("BOT_TOKEN env variable not set")
	}
}
