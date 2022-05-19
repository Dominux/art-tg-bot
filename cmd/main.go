package main

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"

	"artTgBot/internal"
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

	var (
		// Universal markup builders.
		menu = &tele.ReplyMarkup{ResizeKeyboard: true}

		// Reply buttons.
		btnHelp     = menu.Text("ℹ Help")
		btnSettings = menu.Text("⚙ Settings")
	)

	// Filling the keyboard
	menu.Reply(
		menu.Row(btnHelp),
		menu.Row(btnSettings),
	)

	// Creating the handler obj
	handler := internal.NewHandler(menu)

	b.Handle("/start", handler.StartHandler)

	// On reply button pressed (message)
	b.Handle(&btnHelp, func(c tele.Context) error {
		return c.Send("Пошел нахуй пидор ебынй")
	})

	b.Start()
}
