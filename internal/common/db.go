package common

import (
	"gorm.io/gorm"

	"artTgBot/internal/apps/orders"
)

func InitDB(dialector gorm.Dialector, opts ...gorm.Option) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, opts...)
	if err != nil {
		return nil, err
	}

	// Migrations
	db.AutoMigrate(
		// Add here all the models
		&orders.Order{},
	)

	return db, nil
}
