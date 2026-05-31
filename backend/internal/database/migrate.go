package database

import (
	"mini-12306/backend/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := preparePassengerProfileIndexes(db); err != nil {
		return err
	}

	return db.AutoMigrate(
		&model.User{},
		&model.PassengerProfile{},
		&model.Station{},
		&model.Train{},
		&model.TrainStop{},
		&model.Inventory{},
		&model.Order{},
		&model.OrderItem{},
		&model.Payment{},
		&model.Ticket{},
		&model.Refund{},
		&model.ChangeRecord{},
		&model.SystemSetting{},
	)
}

func preparePassengerProfileIndexes(db *gorm.DB) error {
	if db.Dialector.Name() != "mysql" || !db.Migrator().HasTable(&model.PassengerProfile{}) {
		return nil
	}

	const tableName = "passenger_profiles"
	const indexName = "idx_passenger_profiles_user_id"

	var indexCount int64
	if err := db.Raw(`
		SELECT COUNT(*)
		FROM information_schema.statistics
		WHERE table_schema = DATABASE()
		  AND table_name = ?
		  AND index_name = ?
		  AND non_unique = 0
	`, tableName, indexName).Scan(&indexCount).Error; err != nil {
		return err
	}
	if indexCount == 0 {
		return nil
	}

	var constraintNames []string
	if err := db.Raw(`
		SELECT DISTINCT constraint_name
		FROM information_schema.key_column_usage
		WHERE table_schema = DATABASE()
		  AND table_name = ?
		  AND column_name = 'user_id'
		  AND referenced_table_name IS NOT NULL
	`, tableName).Scan(&constraintNames).Error; err != nil {
		return err
	}

	for _, constraintName := range constraintNames {
		if constraintName == "" {
			continue
		}
		if err := db.Exec("ALTER TABLE `passenger_profiles` DROP FOREIGN KEY `" + constraintName + "`").Error; err != nil {
			return err
		}
	}

	if err := db.Exec("DROP INDEX `" + indexName + "` ON `passenger_profiles`").Error; err != nil {
		return err
	}
	return db.Exec("CREATE INDEX `" + indexName + "` ON `passenger_profiles` (`user_id`)").Error
}
