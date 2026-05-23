package database

import (
	"mini-12306/backend/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.PassengerProfile{},
		&model.Station{},
		&model.Train{},
		&model.TrainStop{},
		&model.Inventory{},
		&model.Order{},
		&model.Payment{},
		&model.Ticket{},
	)
}
