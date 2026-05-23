package database

import (
	"time"

	"mini-12306/backend/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type stationSeed struct {
	Code string
	Name string
	City string
}

type stopSeed struct {
	StationCode string
	Order       int
	ArriveClock string
	DepartClock string
	Mileage     int
}

func SeedDemoData(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := seedUsers(tx); err != nil {
			return err
		}
		stations, err := seedStations(tx)
		if err != nil {
			return err
		}
		if err := seedTrainWithInventory(tx, stations); err != nil {
			return err
		}
		return nil
	})
}

func seedUsers(tx *gorm.DB) error {
	var count int64
	if err := tx.Model(&model.User{}).Where("username = ?", "admin").Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("Admin123456"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if err := tx.Create(&model.User{
			Username:     "admin",
			PasswordHash: string(hash),
			Role:         model.UserRoleAdmin,
			Status:       model.UserStatusActive,
		}).Error; err != nil {
			return err
		}
	}

	if err := tx.Model(&model.User{}).Where("username = ?", "alice").Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := model.User{
		Username:     "alice",
		PasswordHash: string(hash),
		Role:         model.UserRolePassenger,
		Status:       model.UserStatusActive,
	}
	if err := tx.Create(&user).Error; err != nil {
		return err
	}
	return tx.Create(&model.PassengerProfile{
		UserID:         user.ID,
		RealName:       "张三",
		IDCardNo:       "110101199001011234",
		Phone:          "13800138000",
		BankCardNo:     "6222020202020202020",
		VerifiedStatus: model.VerificationStatusVerified,
	}).Error
}

func seedStations(tx *gorm.DB) (map[string]model.Station, error) {
	seeds := []stationSeed{
		{Code: "BJN", Name: "北京南", City: "北京"},
		{Code: "TJN", Name: "天津南", City: "天津"},
		{Code: "JNX", Name: "济南西", City: "济南"},
		{Code: "NJH", Name: "南京南", City: "南京"},
		{Code: "SHH", Name: "上海虹桥", City: "上海"},
		{Code: "HZD", Name: "杭州东", City: "杭州"},
	}

	result := make(map[string]model.Station, len(seeds))
	for _, seed := range seeds {
		var station model.Station
		err := tx.Where("code = ?", seed.Code).First(&station).Error
		if err == nil {
			result[seed.Code] = station
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		station = model.Station{
			Code:   seed.Code,
			Name:   seed.Name,
			City:   seed.City,
			Status: model.StationStatusActive,
		}
		if err := tx.Create(&station).Error; err != nil {
			return nil, err
		}
		result[seed.Code] = station
	}
	return result, nil
}
//分支aaa
func seedTrainWithInventory(tx *gorm.DB, stations map[string]model.Station) error {
	train, err := findOrCreateTrain(tx, "G101", "G")
	if err != nil {
		return err
	}
	if err := seedStops(tx, train.ID, stations, []stopSeed{
		{StationCode: "BJN", Order: 1, ArriveClock: "", DepartClock: "08:00:00", Mileage: 0},
		{StationCode: "TJN", Order: 2, ArriveClock: "08:32:00", DepartClock: "08:35:00", Mileage: 122},
		{StationCode: "JNX", Order: 3, ArriveClock: "09:48:00", DepartClock: "09:52:00", Mileage: 406},
		{StationCode: "NJH", Order: 4, ArriveClock: "11:52:00", DepartClock: "11:56:00", Mileage: 1023},
		{StationCode: "SHH", Order: 5, ArriveClock: "13:30:00", DepartClock: "", Mileage: 1318},
	}); err != nil {
		return err
	}
	if err := seedInventories(tx, train.ID, stations["BJN"].ID, stations["SHH"].ID, []inventorySeed{
		{SeatClassCode: "SECOND", PriceCents: 55300, TotalCount: 80, AvailableCount: 32},
		{SeatClassCode: "FIRST", PriceCents: 93300, TotalCount: 32, AvailableCount: 10},
		{SeatClassCode: "BUSINESS", PriceCents: 174800, TotalCount: 8, AvailableCount: 3},
	}); err != nil {
		return err
	}

	train2, err := findOrCreateTrain(tx, "G137", "G")
	if err != nil {
		return err
	}
	if err := seedStops(tx, train2.ID, stations, []stopSeed{
		{StationCode: "BJN", Order: 1, ArriveClock: "", DepartClock: "14:10:00", Mileage: 0},
		{StationCode: "JNX", Order: 2, ArriveClock: "15:58:00", DepartClock: "16:03:00", Mileage: 406},
		{StationCode: "NJH", Order: 3, ArriveClock: "18:06:00", DepartClock: "18:11:00", Mileage: 1023},
		{StationCode: "SHH", Order: 4, ArriveClock: "19:38:00", DepartClock: "19:45:00", Mileage: 1318},
		{StationCode: "HZD", Order: 5, ArriveClock: "20:36:00", DepartClock: "", Mileage: 1477},
	}); err != nil {
		return err
	}
	return seedInventories(tx, train2.ID, stations["BJN"].ID, stations["HZD"].ID, []inventorySeed{
		{SeatClassCode: "SECOND", PriceCents: 62400, TotalCount: 96, AvailableCount: 46},
		{SeatClassCode: "FIRST", PriceCents: 101800, TotalCount: 42, AvailableCount: 18},
	})
}

func findOrCreateTrain(tx *gorm.DB, trainNo string, trainType string) (model.Train, error) {
	var train model.Train
	err := tx.Where("train_no = ?", trainNo).First(&train).Error
	if err == nil {
		return train, nil
	}
	if err != gorm.ErrRecordNotFound {
		return model.Train{}, err
	}

	train = model.Train{
		TrainNo:   trainNo,
		TrainType: trainType,
		Status:    model.TrainStatusActive,
	}
	return train, tx.Create(&train).Error
}

func seedStops(tx *gorm.DB, trainID uint64, stations map[string]model.Station, seeds []stopSeed) error {
	for _, seed := range seeds {
		station := stations[seed.StationCode]
		var count int64
		if err := tx.Model(&model.TrainStop{}).Where("train_id = ? AND station_id = ?", trainID, station.ID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		if err := tx.Create(&model.TrainStop{
			TrainID:     trainID,
			StationID:   station.ID,
			StopOrder:   seed.Order,
			ArriveClock: seed.ArriveClock,
			DepartClock: seed.DepartClock,
			Mileage:     seed.Mileage,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

type inventorySeed struct {
	SeatClassCode  string
	PriceCents     int64
	TotalCount     int
	AvailableCount int
}

func seedInventories(tx *gorm.DB, trainID, fromStationID, toStationID uint64, seeds []inventorySeed) error {
	start, _ := time.ParseInLocation("2006-01-02", time.Now().AddDate(0, 0, 1).Format("2006-01-02"), time.Local)
	for day := 0; day < 7; day++ {
		travelDate := start.AddDate(0, 0, day)
		travelDateText := travelDate.Format("2006-01-02")
		for _, seed := range seeds {
			var count int64
			err := tx.Model(&model.Inventory{}).
				Where("train_id = ? AND DATE(travel_date) = ? AND from_station_id = ? AND to_station_id = ? AND seat_class_code = ?", trainID, travelDateText, fromStationID, toStationID, seed.SeatClassCode).
				Count(&count).Error
			if err != nil {
				return err
			}
			if count > 0 {
				continue
			}

			if err := tx.Create(&model.Inventory{
				TrainID:        trainID,
				TravelDate:     travelDate,
				FromStationID:  fromStationID,
				ToStationID:    toStationID,
				SeatClassCode:  seed.SeatClassCode,
				PriceCents:     seed.PriceCents,
				TotalCount:     seed.TotalCount,
				AvailableCount: seed.AvailableCount,
				LockedCount:    0,
				SoldCount:      seed.TotalCount - seed.AvailableCount,
				Status:         model.InventoryStatusActive,
			}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
