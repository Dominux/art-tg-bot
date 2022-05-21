package orders

import (
	tele "gopkg.in/telebot.v3"

	"artTgBot/internal/apps/info"
	"artTgBot/internal/common"
)

type OrderHandler struct {
	admins []*common.Admin
	menu   *tele.ReplyMarkup
}

func NewHandler(b *tele.Bot, admins []*common.Admin) *OrderHandler {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

	btnCancel := menu.Text("Отмена")
	menu.Reply(menu.Row(btnCancel))

	handler := &OrderHandler{menu: menu, admins: admins}
	b.Handle(&btnCancel, handler.CreateOrder)

	return handler
}

func (h *OrderHandler) Menu(c tele.Context, text string) error {
	return c.Send(text, h.menu)
}

// Order creation
func (h *OrderHandler) CreateOrder(c tele.Context) error {
	// Creating order
	// In current app version we just forward message to admin
	c.ForwardTo(h.admins[0])

	// Getting user back to main info menu
	infoHandler := info.NewHandler(c.Bot())
	return infoHandler.Menu(c, "Спасибо за заказ")
}
