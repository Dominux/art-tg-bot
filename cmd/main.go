package main

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"

	artInfo "artTgBot/internal/apps/info"
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

	// Orders handling

	b.Start()
}
