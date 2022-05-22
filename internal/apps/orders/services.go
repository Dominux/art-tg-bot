package orders

import "gorm.io/gorm"

type OrdersService struct {
	repo *OrdersRepository
}

func NewOrdersService(db *gorm.DB) *OrdersService {
	return &OrdersService{repo: newOrdersRepository(db)}
}

func (o OrdersService) SaveOrder(chatID int64) error {
	return o.repo.createOrder(chatID)
}

func (o OrdersService) GetOrder(chatID int64) (*Order, error) {
	return o.repo.getOrder(chatID)
}

// Submitting order
func (o OrdersService) SubmitOrder(chatID int64) error {
	// In current version it means deleting order from db
	return o.repo.deleteOrder(chatID)
}
