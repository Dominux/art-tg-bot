package info

import (
	"time"

	tele "gopkg.in/telebot.v3"

	"artTgBot/internal/apps/orders"
	"artTgBot/internal/common"
)

type InfoHandler struct {
	menu *tele.ReplyMarkup
}

func NewHandler(b *tele.Bot) *InfoHandler {
	var (
		// Universal markup builders.
		menu = &tele.ReplyMarkup{ResizeKeyboard: true}

		// Reply buttons.
		btnShowExamples = menu.Text("ℹ Примеры paбот")
		btnCreateOrder  = menu.Text("Заказать работу")
		// btnSettings = menu.Text("⚙ Settings")
	)

	// Filling the keyboard
	menu.Reply(
		menu.Row(btnShowExamples),
		menu.Row(btnCreateOrder),
	)

	handler := &InfoHandler{menu: menu}

	b.Handle(&btnShowExamples, handler.showExamples)
	b.Handle(&btnCreateOrder, handler.showExamples)

	return handler
}

func (h InfoHandler) HandleStart(c tele.Context) error {
	return c.Send("Добро пожаловать!", h.menu)
}

func (h InfoHandler) showExamples(c tele.Context) error {
	if err := c.Send("Lol"); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)
	return c.Send("Какая-то ссылка на арты")
}

func (h InfoHandler) createOrder(c tele.Context) error {
	orderHandler := orders.NewHandler(c.Bot(), []*common.Admin{})
	return orderHandler.Menu(
		c,
		"Опишите ваш заказ как можно подробнее\nПо возможности приложите к сообщению ссылку или изображение того, из чего хотите получить Арт",
	)
}

func (h *InfoHandler) Menu(c tele.Context, text string) error {
	return c.Send(text, h.menu)
}
