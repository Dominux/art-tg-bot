package orders

import "time"

type Order struct {
	ChatID      int64 `gorm:"primarykey"`
	CreatedAt   time.Time
	MessagesIDs string
}
