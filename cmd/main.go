package main

import (
	"fmt"
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
		btnSubmitOrder := orderMenu.Text("Отправить заказ")
		btnCancel := orderMenu.Text("Отмена")
		orderMenu.Reply(
			orderMenu.Row(btnSubmitOrder),
			orderMenu.Row(btnCancel),
		)

		b.Handle(&btnSubmitOrder, submitOrder)
		b.Handle(&btnCancel, cancelOrderCreation)
	}

	// Handling any message content here
	{
		hooks := []string{
			tele.OnText,
			tele.OnPhoto,
			tele.OnAudio,
			tele.OnAnimation,
			tele.OnDocument,
			tele.OnSticker,
			tele.OnVideo,
			tele.OnVoice,
			tele.OnVideoNote,
			tele.OnContact,
			tele.OnLocation,
			tele.OnVenue,
			tele.OnDice,
			tele.OnInvoice,
			tele.OnPayment,
			tele.OnGame,
			tele.OnPoll,
			tele.OnPollAnswer,
		}
		for _, hook := range hooks {
			b.Handle(hook, onMessage)
		}
	}

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
	if err := ordersService.CreateOrder(chatID); err != nil {
		return c.Send("Ошибка при создании заказа, повторите попытку позже", infoMenu)
	}

	return c.Send(
		`Опишите ваш заказ как можно подробнее
По возможности приложите к сообщению ссылку или изображение того, из чего хотите получить Арт`,
		orderMenu,
	)
}

func onMessage(c tele.Context) error {
	chatID := c.Chat().ID
	messageID := c.Message().ID

	// Trying to add detail to order, if order isn't found - then it's just a user mistake
	if err := ordersService.AddDetail(chatID, messageID); err != nil {
		return c.Send(
			`Не понял вас :)
Пожалуйста, используйте клавиатуру для навигации`,
			infoMenu,
		)
	}

	return nil
}

func submitOrder(c tele.Context) error {
	chatID := c.Chat().ID

	errMsg := "Произошла ошибка с обработкой заказа, повторите позже"

	// In current app version we just forward messages to admin
	messagesIDs, err := ordersService.GetOrdersDetails(chatID)
	if err != nil {
		return c.Send(errMsg, infoMenu)
	}

	if err := func() error {
		// TODO: handle error somehow
		defer ordersService.DeleteOrder(chatID)

		// Notifying admin about new order
		_, err := c.Bot().Send(admin, fmt.Sprintf("Новый заказ от пользователя @%s", c.Chat().Username))
		if err != nil {
			fmt.Println(err)
			return err
		}

		for _, messageID := range messagesIDs {
			storedMessage := tele.StoredMessage{MessageID: messageID, ChatID: chatID}
			_, err := c.Bot().Forward(admin, storedMessage)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}

		return nil
	}(); err != nil {
		return c.Send(errMsg, infoMenu)
	}

	// Getting user back to main info menu
	return c.Send("Спасибо за заказ", infoMenu)
}

func cancelOrderCreation(c tele.Context) error {
	chatID := c.Chat().ID
	// TODO: handle error
	ordersService.DeleteOrder(chatID)

	return c.Send("Вы вернулись в основное меню", infoMenu)
}
