package internal

import (
	tele "gopkg.in/telebot.v3"
)

type Handler struct {
	menu *tele.ReplyMarkup
}

func NewHandler(menu *tele.ReplyMarkup) *Handler {
	return &Handler{menu: menu}
}

func (h *Handler) StartHandler(c tele.Context) error {
	return c.Send("Добро пожаловать!", h.menu)
}

func (h *Handler) ShowArts(c tele.Context) error {
	return c.Send("Какая-то ссылка на арты")
}
