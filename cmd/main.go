package main

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"

	artInfo "artTgBot/internal/info"
)

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Info handling
	infoHandler := artInfo.NewHandler(b)
	b.Handle("/start", infoHandler.HandleStart)

	// Orders handling

	b.Start()
}
