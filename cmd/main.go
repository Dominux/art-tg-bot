package main

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"artTgBot/internal/apps/orders"
	"artTgBot/internal/common"
)

var (
	ordersService *orders.OrdersService

	admin *common.Admin

	infoMenu  = &tele.ReplyMarkup{ResizeKeyboard: true}
	orderMenu = &tele.ReplyMarkup{ResizeKeyboard: true}
)

func main() {
	// Creating bot
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

	// Initiating db and services
	{
		dbName := os.Getenv("DB_NAME")
		db, err := common.InitDB(sqlite.Open(dbName), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		ordersService = orders.NewOrdersService(db)
	}

	// Creating admins
	{
		adminStr := os.Getenv("ADMIN")
		admin = common.NewAdmin(adminStr)
	}

	// Info handling
	{
		var (
			// Reply buttons.
			btnShowExamples = infoMenu.Text("ℹ Примеры paбот")
			btnCreateOrder  = infoMenu.Text("Заказать работу")
		)

		// Filling the keyboard
		infoMenu.Reply(
			infoMenu.Row(btnShowExamples),
			infoMenu.Row(btnCreateOrder),
		)

		b.Handle(&btnShowExamples, showExamples)
		b.Handle(&btnCreateOrder, getOrderCreationForm)
	}

	b.Handle("/start", handleStart)

	// Orders handling
	{
		btnCancel := orderMenu.Text("Отмена")
		orderMenu.Reply(orderMenu.Row(btnCancel))

		b.Handle(&btnCancel, cancelOrderCreation)
	}

	// Handling any text
	// TODO: bug with handling only text messages
	b.Handle(tele.OnText, onText)

	// Starting bot
	b.Start()
}

//////////////////////////////////////////////////////////////////////////////
// 	Handlers
//////////////////////////////////////////////////////////////////////////////

func handleStart(c tele.Context) error {
	return c.Send("Добро пожаловать!", infoMenu)
}

func showExamples(c tele.Context) error {
	if err := c.Send("Lol"); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)
	return c.Send("Какая-то ссылка на арты")
}

func getOrderCreationForm(c tele.Context) error {
	chatID := c.Chat().ID
	if err := ordersService.SaveOrder(chatID); err != nil {
		return c.Send("Ошибка при создании заказа, повторите попытку позже", infoMenu)
	}

	return c.Send(
		`Опишите ваш заказ как можно подробнее
		По возможности приложите к сообщению ссылку или изображение того, из чего хотите получить Арт`,
		orderMenu,
	)
}

func onText(c tele.Context) error {
	chatID := c.Chat().ID

	// Trying to get order, if order isn't found - then it's just a user mistake
	_, err := ordersService.GetOrder(chatID)
	if err != nil {
		return c.Send(
			`Не понял вас :)
			Пожалуйста, используйте клавиатуру для навигации`,
			infoMenu,
		)
	}

	// Otherwise - creating order
	return createOrder(c, chatID)
}

func createOrder(c tele.Context, chatID int64) error {
	// Creating order
	// In current app version we just forward message to admin
	err := c.ForwardTo(admin)

	// Submitting order
	// TODO: handle error somehow
	ordersService.SubmitOrder(chatID)

	// Handling error after submiting order cause otherwise it may cause some wrong behavior
	if err != nil {
		return c.Send("Произошла ошибка с обработкой заказа, повторите позже", infoMenu)
	}

	// Getting user back to main info menu
	return c.Send("Спасибо за заказ", infoMenu)
}

func cancelOrderCreation(c tele.Context) error {
	return c.Send("Вы вернулись в основное меню", infoMenu)
}
