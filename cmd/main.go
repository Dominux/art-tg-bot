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

	mainMenu                  = &tele.ReplyMarkup{ResizeKeyboard: true}
	portfolioMenu             = &tele.ReplyMarkup{ResizeKeyboard: true}
	portfolioByCategoriesMenu = &tele.ReplyMarkup{ResizeKeyboard: true}
	orderArtMenu              = &tele.ReplyMarkup{ResizeKeyboard: true}
	orderArtTypeMenu          = &tele.ReplyMarkup{ResizeKeyboard: true}
	submitOrderMenu           = &tele.ReplyMarkup{ResizeKeyboard: true}
)

func main() {
	// Creating bot
	var b *tele.Bot
	var err error
	// switch os.Getenv("MODE") {
	// case "PRODUCTION":
	// 	panic("Not implemented")

	// default:
	pref := tele.Settings{
		Token:  os.Getenv("API_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err = tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	// }

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

	b.Handle("/start", handleStart)

	// Info handling
	{
		var (
			// Reply buttons.
			btnAboutProject  = mainMenu.Text("О проекте")
			btnActions       = mainMenu.Text("Акции")
			btnPrices        = mainMenu.Text("Цены")
			btnShowReviews   = mainMenu.Text("Отзывы")
			btnShowPortfolio = mainMenu.Text("Посмотреть портфолио")
			btnOrderArt      = mainMenu.Text("Заказать Арт")
		)

		// Filling the keyboard
		mainMenu.Reply(
			mainMenu.Row(btnAboutProject),
			mainMenu.Row(btnActions),
			mainMenu.Row(btnPrices),
			mainMenu.Row(btnShowReviews),
			mainMenu.Row(btnShowPortfolio),
			mainMenu.Row(btnOrderArt),
		)

		b.Handle(&btnAboutProject, aboutProjects)
		b.Handle(&btnActions, showActions)
		b.Handle(&btnPrices, showPrices)
		b.Handle(&btnShowReviews, showReviews)
		b.Handle(&btnShowPortfolio, showPortfolio)
		b.Handle(&btnOrderArt, orderArt)
	}

	// Portfolio handling
	{
		btnRandomArt := portfolioMenu.Text("Случайный Арт")
		btnPortfolioByCategories := portfolioMenu.Text("Портфолио по категориям")
		btnBackToMainMenu := portfolioMenu.Text("Назад")

		portfolioMenu.Reply(
			portfolioMenu.Row(btnPortfolioByCategories),
			portfolioMenu.Row(btnRandomArt),
			portfolioMenu.Row(btnBackToMainMenu),
		)

		b.Handle(&btnPortfolioByCategories, showPortfolioByCategories)
		b.Handle(&btnRandomArt, showRandomArt)
		b.Handle(&btnBackToMainMenu, backToMainMenu)
	}

	// Portfolio by categories handling
	{
		btnUltraArt := portfolioByCategoriesMenu.Text("Арт Ультра")
		btnStickers := portfolioByCategoriesMenu.Text("Стикеры")
		btnStandardPicture := portfolioByCategoriesMenu.Text("Стандартный рисунок")
		btnThreeD := portfolioByCategoriesMenu.Text("3D")
		btnEZArt := portfolioByCategoriesMenu.Text("EZ Арт")
		btnBanners := portfolioByCategoriesMenu.Text("Баннеры")
		btnBackToMainMenu := portfolioByCategoriesMenu.Text("Назад")

		portfolioByCategoriesMenu.Reply(
			portfolioByCategoriesMenu.Row(btnUltraArt),
			portfolioByCategoriesMenu.Row(btnStickers),
			portfolioByCategoriesMenu.Row(btnStandardPicture),
			portfolioByCategoriesMenu.Row(btnThreeD),
			portfolioByCategoriesMenu.Row(btnEZArt),
			portfolioByCategoriesMenu.Row(btnBanners),
			portfolioByCategoriesMenu.Row(btnBackToMainMenu),
		)

		b.Handle(&btnUltraArt, artUltraCategory)
		b.Handle(&btnStickers, stickersCategory)
		b.Handle(&btnStandardPicture, standardPictureCategory)
		b.Handle(&btnThreeD, threeDCategory)
		b.Handle(&btnEZArt, ezArtCategory)
		b.Handle(&btnBanners, bannersCategory)
		b.Handle(&btnBackToMainMenu, backToMainMenu)
	}

	// Order art handling
	{
		btnConditions := orderArtMenu.Text("Условия")
		btnChooseArtType := orderArtMenu.Text("Выбрать тип Арта")
		btnBackToMainMenu := orderArtMenu.Text("Назад")

		orderArtMenu.Reply(
			orderArtMenu.Row(btnConditions),
			orderArtMenu.Row(btnChooseArtType),
			orderArtMenu.Row(btnBackToMainMenu),
		)

		b.Handle(&btnConditions, showConditions)
		b.Handle(&btnChooseArtType, chooseArtType)
		b.Handle(&btnBackToMainMenu, backToMainMenu)
	}

	// Choose art type handling
	{
		btnUltraArt := orderArtTypeMenu.Text("Категория Арт ультра")
		btnStickers := orderArtTypeMenu.Text("Категория Стикеры")
		btnStandardPicture := orderArtTypeMenu.Text("Категория Стандартный рисунок")
		btnThreeD := orderArtTypeMenu.Text("Категория 3D")
		btnEZArt := orderArtTypeMenu.Text("Категория EZ Арт")
		btnBanners := orderArtTypeMenu.Text("Категория Баннеры")
		btnBackToMainMenu := orderArtTypeMenu.Text("Назад")

		orderArtTypeMenu.Reply(
			orderArtTypeMenu.Row(btnUltraArt),
			orderArtTypeMenu.Row(btnStickers),
			orderArtTypeMenu.Row(btnStandardPicture),
			orderArtTypeMenu.Row(btnThreeD),
			orderArtTypeMenu.Row(btnEZArt),
			orderArtTypeMenu.Row(btnBanners),
			orderArtTypeMenu.Row(btnBackToMainMenu),
		)

		b.Handle(&btnUltraArt, getOrderCreationForm)
		b.Handle(&btnStickers, getOrderCreationForm)
		b.Handle(&btnStandardPicture, getOrderCreationForm)
		b.Handle(&btnThreeD, getOrderCreationForm)
		b.Handle(&btnEZArt, getOrderCreationForm)
		b.Handle(&btnBanners, getOrderCreationForm)
		b.Handle(&btnBackToMainMenu, backToMainMenu)
	}

	// Submitting order handling
	{
		btnSubmitOrder := submitOrderMenu.Text("Создать заказ")
		btnCancelOrder := submitOrderMenu.Text("Отменить")

		submitOrderMenu.Reply(
			submitOrderMenu.Row(btnSubmitOrder),
			submitOrderMenu.Row(btnCancelOrder),
		)

		b.Handle(&btnSubmitOrder, submitOrder)
		b.Handle(&btnCancelOrder, cancelOrderCreation)
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
		}
		for _, hook := range hooks {
			b.Handle(hook, onContentMessage)
		}
	}

	// Starting bot
	b.Start()
}

//////////////////////////////////////////////////////////////////////////////
// 	Menu handlers
//////////////////////////////////////////////////////////////////////////////

func handleStart(c tele.Context) error {
	return c.Send("Добро пожаловать!", mainMenu)
}

func backToMainMenu(c tele.Context) error {
	return c.Send("Вы вернулись в основное меню", mainMenu)
}

func aboutProjects(c tele.Context) error {
	return c.Send("Топ проджект ин да уорлд!")
}

func showActions(c tele.Context) error {
	return c.Send("На данный момент акций нет")
}

func showPrices(c tele.Context) error {
	return c.Send("300$")
}

func showReviews(c tele.Context) error {
	return c.Send("Отзывы")
}

func showPortfolio(c tele.Context) error {
	return c.Send("Наше портфолио", portfolioMenu)
}

func orderArt(c tele.Context) error {
	return c.Send("Здесь вы можете заказать Арт", orderArtMenu)
}

//////////////////////////////////////////////////////////////////////////////
// 	Portfolio handlers
//////////////////////////////////////////////////////////////////////////////

func showRandomArt(c tele.Context) error {
	return c.Send("Рандомный арт")
}

func showPortfolioByCategories(c tele.Context) error {
	return c.Send("Портфолио по категориям", portfolioByCategoriesMenu)
}

//////////////////////////////////////////////////////////////////////////////
// 	Portfolio by categories handlers
//////////////////////////////////////////////////////////////////////////////

func artUltraCategory(c tele.Context) error {
	return c.Send("Арт ультра работы")
}

func stickersCategory(c tele.Context) error {
	return c.Send("kjk")
}

func standardPictureCategory(c tele.Context) error {
	return c.Send("Стандартный рисунок")
}

func threeDCategory(c tele.Context) error {
	return c.Send("Арт ультра работы")
}

func ezArtCategory(c tele.Context) error {
	return c.Send("Арт ультра работы")
}

func bannersCategory(c tele.Context) error {
	return c.Send("Арт ультра работы")
}

//////////////////////////////////////////////////////////////////////////////
// 	Art ordering handlers
//////////////////////////////////////////////////////////////////////////////

func showConditions(c tele.Context) error {
	return c.Send("Рандомный арт")
}

func chooseArtType(c tele.Context) error {
	return c.Send("Портфолио по категориям", orderArtTypeMenu)
}

//////////////////////////////////////////////////////////////////////////////
// 	Order handlers
//////////////////////////////////////////////////////////////////////////////

func getOrderCreationForm(c tele.Context) error {
	chatID := c.Chat().ID
	messageID := c.Message().ID

	// Creating order
	if err := ordersService.CreateOrder(chatID); err != nil {
		return c.Send("Ошибка при создании заказа, повторите попытку позже", orderArtTypeMenu)
	}

	// Saving category the user chose
	if err := ordersService.AddDetail(chatID, messageID); err != nil {
		return c.Send("Ошибка при создании заказа, повторите попытку позже", orderArtTypeMenu)
	}

	return c.Send(
		`Опишите ваш заказ как можно подробнее
По возможности приложите к сообщению ссылку или изображение того, из чего хотите получить Арт`,
		submitOrderMenu,
	)
}

func onContentMessage(c tele.Context) error {
	chatID := c.Chat().ID
	messageID := c.Message().ID

	// Trying to add detail to order, if order isn't found - then it's just a user mistake
	if err := ordersService.AddDetail(chatID, messageID); err != nil {
		return c.Send(
			`Не понял вас :)
Пожалуйста, используйте клавиатуру для навигации`,
			orderArtTypeMenu,
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
		return c.Send(errMsg, mainMenu)
	}

	// If there's no messages
	if len(messagesIDs) == 0 {
		return c.Send("Вы не можете создать пустой заказ")
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
		return c.Send(errMsg, mainMenu)
	}

	// Getting user back to main info menu
	return c.Send("Спасибо за заказ", mainMenu)
}

func cancelOrderCreation(c tele.Context) error {
	chatID := c.Chat().ID
	// TODO: handle error
	ordersService.DeleteOrder(chatID)

	return c.Send("Вы вернулись в основное меню", mainMenu)
}
