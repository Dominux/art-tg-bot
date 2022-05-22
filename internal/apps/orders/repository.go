package orders

import (
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
	order := &Order{ChatID: chatID, CreatedAt: time.Now()}
	return o.db.Create(order).Error
}

func (o OrdersRepository) deleteOrder(chatID int64) error {
	return o.db.Delete(&Order{}, chatID).Error
}
