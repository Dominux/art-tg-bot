package main

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"

	artInfo "artTgBot/internal/apps/info"
	artOrders "artTgBot/internal/apps/orders"
	"artTgBot/internal/common"
)

func main() {

	var b *tele.Bot
	var err error
	switch os.Getenv("MODE") {
	case "PRODUCTION":
		panic("Not implemented")

	default:
		pref := tele.Settings{
			Token:  os.Getenv("API_TOKEN"),
			Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		}

		b, err = tele.NewBot(pref)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	// Info handling
	infoHandler := artInfo.NewHandler(b)
	b.Handle("/start", infoHandler.HandleStart)

	// Creating admins
	adminStr := os.Getenv("ADMIN")
	admin := common.NewAdmin(adminStr)

	// Orders handling
	orderHandler := artOrders.NewHandler(b, []*common.Admin{admin})

	b.Start()
}
