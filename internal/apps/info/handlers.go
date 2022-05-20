package info

import (
	tele "gopkg.in/telebot.v3"
)

type Handler struct {
	menu *tele.ReplyMarkup
}

func NewHandler(b *tele.Bot) *Handler {
	var (
		// Universal markup builders.
		menu = &tele.ReplyMarkup{ResizeKeyboard: true}

		// Reply buttons.
		btnShowExamples = menu.Text("ℹ Примеры paбот")
		// btnSettings = menu.Text("⚙ Settings")
	)

	// Filling the keyboard
	menu.Reply(
		menu.Row(btnShowExamples),
		// menu.Row(btnSettings),
	)

	handler := &Handler{menu: menu}

	b.Handle(&btnShowExamples, handler.showExamples)

	return handler
}

func (h *Handler) HandleStart(c tele.Context) error {
	return c.Send("Добро пожаловать!", h.menu)
}

func (h *Handler) showExamples(c tele.Context) error {
	return c.Send("Какая-то ссылка на арты")
}
