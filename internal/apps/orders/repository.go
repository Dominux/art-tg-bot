package orders

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type OrdersRepository struct {
	db *gorm.DB
}

func newOrdersRepository(db *gorm.DB) *OrdersRepository {
	return &OrdersRepository{db: db}
}

func (o OrdersRepository) createOrder(chatID int64) error {
	order := &Order{ChatID: chatID, CreatedAt: time.Now(), MessagesIDs: ""}
	return o.db.Create(order).Error
}

func (o OrdersRepository) getOrder(chatID int64) (*Order, error) {
	order := &Order{}
	if err := o.db.First(order, chatID).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (o OrdersRepository) addDetail(chatID int64, messageID int) error {
	// Getting order
	order, err := o.getOrder(chatID)
	if err != nil {
		return err
	}

	// Updating details
	if len(order.MessagesIDs) > 0 {
		order.MessagesIDs += "," + strconv.Itoa(messageID)
	} else {
		order.MessagesIDs = strconv.Itoa(messageID)
	}

	return o.db.Model(order).Update("MessagesIDs", order.MessagesIDs).Error
}

func (o OrdersRepository) deleteOrder(chatID int64) error {
	return o.db.Delete(&Order{}, chatID).Error
}
