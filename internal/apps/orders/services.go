package orders

import (
	"strings"

	"gorm.io/gorm"
)

type OrdersService struct {
	repo *OrdersRepository
}

func NewOrdersService(db *gorm.DB) *OrdersService {
	return &OrdersService{repo: newOrdersRepository(db)}
}

func (o OrdersService) CreateOrder(chatID int64) error {
	return o.repo.createOrder(chatID)
}

func (o OrdersService) AddDetail(chatID int64, messageID int) error {
	return o.repo.addDetail(chatID, messageID)
}

// Submitting order
func (o OrdersService) GetOrdersDetails(chatID int64) ([]string, error) {
	order, err := o.repo.getOrder(chatID)
	if err != nil {
		return nil, err
	}

	messagesIDs := strings.Split(order.MessagesIDs, ",")
	return messagesIDs, nil
}

func (o OrdersService) DeleteOrder(chatID int64) error {
	return o.repo.deleteOrder(chatID)
}
