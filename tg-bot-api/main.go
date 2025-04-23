package tg_bot_api

import "os"

var (
	token string
)

func init() {
	token = os.Getenv("BOT_TOKEN")
}
