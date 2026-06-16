package database

import (
	"fmt"
	"math"
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

type routeSeed struct {
	TrainNo     string
	TrainType   string
	Stops       []stopSeed
	FromStation string
	ToStation   string
	Inventories []inventorySeed
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

	if err := tx.Model(&model.User{}).Where("username = ?", "clerk").Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("Clerk123456"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if err := tx.Create(&model.User{
			Username:     "clerk",
			PasswordHash: string(hash),
			Role:         model.UserRoleClerk,
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
		{Code: "XZE", Name: "徐州东", City: "徐州"},
		{Code: "NJH", Name: "南京南", City: "南京"},
		{Code: "SZB", Name: "苏州北", City: "苏州"},
		{Code: "SHH", Name: "上海虹桥", City: "上海"},
		{Code: "HZD", Name: "杭州东", City: "杭州"},
		{Code: "NGH", Name: "宁波", City: "宁波"},
		{Code: "HFN", Name: "合肥南", City: "合肥"},
		{Code: "WHN", Name: "武汉", City: "武汉"},
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

func seedTrainWithInventory(tx *gorm.DB, stations map[string]model.Station) error {
	for _, seed := range demoTrainSeeds() {
		train, err := findOrCreateTrain(tx, seed.TrainNo, seed.TrainType)
		if err != nil {
			return err
		}
		if err := seedStops(tx, train.ID, stations, seed.Stops); err != nil {
			return err
		}
		if err := seedRouteInventories(tx, train.ID, stations, seed.Stops, seed.Inventories); err != nil {
			return err
		}
	}
	return nil
}

func demoTrainSeeds() []routeSeed {
	seeds := []routeSeed{
		{
			TrainNo: "G101", TrainType: "G", FromStation: "BJN", ToStation: "SHH",
			Stops: []stopSeed{
				{StationCode: "BJN", Order: 1, ArriveClock: "", DepartClock: "08:00:00", Mileage: 0},
				{StationCode: "TJN", Order: 2, ArriveClock: "08:32:00", DepartClock: "08:35:00", Mileage: 122},
				{StationCode: "JNX", Order: 3, ArriveClock: "09:48:00", DepartClock: "09:52:00", Mileage: 406},
				{StationCode: "XZE", Order: 4, ArriveClock: "10:35:00", DepartClock: "10:40:00", Mileage: 692},
				{StationCode: "NJH", Order: 5, ArriveClock: "11:26:00", DepartClock: "11:30:00", Mileage: 1023},
				{StationCode: "SZB", Order: 6, ArriveClock: "12:03:00", DepartClock: "12:07:00", Mileage: 1192},
				{StationCode: "SHH", Order: 7, ArriveClock: "13:30:00", DepartClock: "", Mileage: 1318},
			},
			Inventories: []inventorySeed{
				{SeatClassCode: "SECOND", TotalCount: 80, AvailableCount: 32},
				{SeatClassCode: "FIRST", TotalCount: 32, AvailableCount: 10},
				{SeatClassCode: "BUSINESS", TotalCount: 8, AvailableCount: 3},
			},
		},
		{
			TrainNo: "G137", TrainType: "G", FromStation: "SHH", ToStation: "HZD",
			Stops: []stopSeed{
				{StationCode: "SHH", Order: 1, ArriveClock: "", DepartClock: "14:10:00", Mileage: 0},
				{StationCode: "SZB", Order: 2, ArriveClock: "14:52:00", DepartClock: "14:55:00", Mileage: 126},
				{StationCode: "NJH", Order: 3, ArriveClock: "15:48:00", DepartClock: "15:52:00", Mileage: 295},
				{StationCode: "XZE", Order: 4, ArriveClock: "16:35:00", DepartClock: "16:39:00", Mileage: 626},
				{StationCode: "JNX", Order: 5, ArriveClock: "17:22:00", DepartClock: "17:26:00", Mileage: 912},
				{StationCode: "TJN", Order: 6, ArriveClock: "18:18:00", DepartClock: "18:21:00", Mileage: 1196},
				{StationCode: "BJN", Order: 7, ArriveClock: "19:46:00", DepartClock: "19:50:00", Mileage: 1318},
				{StationCode: "HZD", Order: 8, ArriveClock: "20:36:00", DepartClock: "", Mileage: 1477},
			},
			Inventories: []inventorySeed{
				{SeatClassCode: "SECOND", TotalCount: 96, AvailableCount: 46},
				{SeatClassCode: "FIRST", TotalCount: 42, AvailableCount: 18},
			},
		},
		{
			TrainNo: "G215", TrainType: "G", FromStation: "BJN", ToStation: "HZD",
			Stops: []stopSeed{
				{StationCode: "BJN", Order: 1, ArriveClock: "", DepartClock: "07:05:00", Mileage: 0},
				{StationCode: "TJN", Order: 2, ArriveClock: "07:38:00", DepartClock: "07:41:00", Mileage: 122},
				{StationCode: "JNX", Order: 3, ArriveClock: "08:36:00", DepartClock: "08:39:00", Mileage: 406},
				{StationCode: "NJH", Order: 4, ArriveClock: "09:24:00", DepartClock: "09:28:00", Mileage: 692},
				{StationCode: "SHH", Order: 5, ArriveClock: "10:32:00", DepartClock: "10:36:00", Mileage: 1318},
				{StationCode: "HZD", Order: 6, ArriveClock: "11:28:00", DepartClock: "", Mileage: 1477},
			},
			Inventories: []inventorySeed{
				{SeatClassCode: "SECOND", TotalCount: 88, AvailableCount: 26},
				{SeatClassCode: "FIRST", TotalCount: 36, AvailableCount: 12},
				{SeatClassCode: "BUSINESS", TotalCount: 10, AvailableCount: 4},
			},
		},
		{
			TrainNo: "D611", TrainType: "D", FromStation: "WHN", ToStation: "SHH",
			Stops: []stopSeed{
				{StationCode: "WHN", Order: 1, ArriveClock: "", DepartClock: "08:15:00", Mileage: 0},
				{StationCode: "HFN", Order: 2, ArriveClock: "09:02:00", DepartClock: "09:05:00", Mileage: 176},
				{StationCode: "NJH", Order: 3, ArriveClock: "10:18:00", DepartClock: "10:22:00", Mileage: 531},
				{StationCode: "SZB", Order: 4, ArriveClock: "11:08:00", DepartClock: "11:12:00", Mileage: 700},
				{StationCode: "SHH", Order: 5, ArriveClock: "12:46:00", DepartClock: "", Mileage: 826},
			},
			Inventories: []inventorySeed{
				{SeatClassCode: "SECOND", TotalCount: 120, AvailableCount: 44},
				{SeatClassCode: "SECOND_SLEEPER", TotalCount: 64, AvailableCount: 20},
				{SeatClassCode: "FIRST_SLEEPER", TotalCount: 32, AvailableCount: 9},
			},
		},
		{
			TrainNo: "K873", TrainType: "K", FromStation: "NJH", ToStation: "WHN",
			Stops: []stopSeed{
				{StationCode: "NJH", Order: 1, ArriveClock: "", DepartClock: "09:35:00", Mileage: 0},
				{StationCode: "HFN", Order: 2, ArriveClock: "10:22:00", DepartClock: "10:25:00", Mileage: 355},
				{StationCode: "WHN", Order: 3, ArriveClock: "12:10:00", DepartClock: "", Mileage: 884},
			},
			Inventories: []inventorySeed{
				{SeatClassCode: "HARD_SEAT", TotalCount: 120, AvailableCount: 50},
				{SeatClassCode: "NO_SEAT", TotalCount: 80, AvailableCount: 28},
				{SeatClassCode: "HARD_SLEEPER", TotalCount: 72, AvailableCount: 18},
				{SeatClassCode: "SOFT_SLEEPER", TotalCount: 32, AvailableCount: 8},
				{SeatClassCode: "DELUXE_SOFT_SLEEPER", TotalCount: 12, AvailableCount: 2},
			},
		},
	}

	patterns := []struct {
		prefix string
		start  int
		count  int
		route  []string
		miles  []int
	}{
		{prefix: "G", start: 201, count: 10, route: []string{"BJN", "TJN", "JNX", "XZE", "NJH", "SZB", "SHH", "HZD"}, miles: []int{0, 122, 406, 692, 1023, 1192, 1318, 1477}},
		{prefix: "G", start: 302, count: 10, route: []string{"SHH", "SZB", "NJH", "XZE", "JNX", "TJN", "BJN"}, miles: []int{0, 126, 295, 626, 912, 1196, 1318}},
		{prefix: "D", start: 501, count: 10, route: []string{"BJN", "JNX", "XZE", "NJH", "HFN", "WHN"}, miles: []int{0, 406, 692, 1023, 1180, 1540}},
		{prefix: "D", start: 602, count: 10, route: []string{"WHN", "HFN", "NJH", "SZB", "SHH", "NGH"}, miles: []int{0, 360, 517, 686, 812, 1126}},
		{prefix: "Z", start: 35, count: 6, route: []string{"BJN", "JNX", "NJH", "SHH", "HZD"}, miles: []int{0, 406, 1023, 1318, 1477}},
		{prefix: "T", start: 109, count: 6, route: []string{"BJN", "TJN", "JNX", "NJH", "SHH"}, miles: []int{0, 122, 406, 1023, 1318}},
		{prefix: "K", start: 701, count: 10, route: []string{"HZD", "SHH", "SZB", "NJH", "HFN", "WHN"}, miles: []int{0, 159, 285, 454, 611, 971}},
	}

	for _, pattern := range patterns {
		for i := 0; i < pattern.count; i++ {
			baseMinutes := 6*60 + i*42
			trainNo := fmt.Sprintf("%s%d", pattern.prefix, pattern.start+i*2)
			seeds = append(seeds, routeSeed{
				TrainNo:     trainNo,
				TrainType:   pattern.prefix,
				Stops:       buildStops(pattern.route, pattern.miles, baseMinutes),
				FromStation: pattern.route[0],
				ToStation:   pattern.route[len(pattern.route)-1],
				Inventories: buildInventories(pattern.prefix, pattern.miles[len(pattern.miles)-1], i),
			})
		}
	}

	return seeds
}

func buildStops(stations []string, miles []int, baseMinutes int) []stopSeed {
	stops := make([]stopSeed, 0, len(stations))
	for i, stationCode := range stations {
		arriveClock := ""
		departClock := ""
		stopMinutes := baseMinutes + i*48 + i*i*4
		if i > 0 {
			arriveClock = clockText(stopMinutes)
		}
		if i < len(stations)-1 {
			departClock = clockText(stopMinutes + 4)
			if i == 0 {
				departClock = clockText(baseMinutes)
			}
		}
		stops = append(stops, stopSeed{
			StationCode: stationCode,
			Order:       i + 1,
			ArriveClock: arriveClock,
			DepartClock: departClock,
			Mileage:     miles[i],
		})
	}
	return stops
}

func clockText(minutes int) string {
	minutes = minutes % (24 * 60)
	return fmt.Sprintf("%02d:%02d:00", minutes/60, minutes%60)
}

func buildInventories(trainType string, _ int, index int) []inventorySeed {
	switch trainType {
	case "D":
		return []inventorySeed{
			{SeatClassCode: "SECOND", TotalCount: 120, AvailableCount: 72 - index%20},
			{SeatClassCode: "SECOND_SLEEPER", TotalCount: 64, AvailableCount: 30 - index%10},
			{SeatClassCode: "FIRST_SLEEPER", TotalCount: 32, AvailableCount: 16 - index%8},
		}
	case "Z", "T", "K":
		return []inventorySeed{
			{SeatClassCode: "HARD_SEAT", TotalCount: 120, AvailableCount: 70 - index%18},
			{SeatClassCode: "NO_SEAT", TotalCount: 80, AvailableCount: 55 - index%15},
			{SeatClassCode: "HARD_SLEEPER", TotalCount: 72, AvailableCount: 36 - index%10},
			{SeatClassCode: "SOFT_SLEEPER", TotalCount: 32, AvailableCount: 14 - index%6},
			{SeatClassCode: "DELUXE_SOFT_SLEEPER", TotalCount: 12, AvailableCount: 5 - index%3},
		}
	default:
		return []inventorySeed{
			{SeatClassCode: "SECOND", TotalCount: 120, AvailableCount: 72 - index%20},
			{SeatClassCode: "FIRST", TotalCount: 48, AvailableCount: 24 - index%12},
			{SeatClassCode: "BUSINESS", TotalCount: 12, AvailableCount: 6 - index%4},
		}
	}
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
	if err := tx.Unscoped().Where("train_id = ?", trainID).Delete(&model.TrainStop{}).Error; err != nil {
		return err
	}
	for _, seed := range seeds {
		station := stations[seed.StationCode]
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
	TotalCount     int
	AvailableCount int
}

func seedRouteInventories(tx *gorm.DB, trainID uint64, stations map[string]model.Station, stops []stopSeed, seeds []inventorySeed) error {
	var train model.Train
	if err := tx.First(&train, trainID).Error; err != nil {
		return err
	}

	start, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
	end := start.AddDate(0, 0, 7)
	if err := tx.Unscoped().
		Where("train_id = ? AND travel_date >= ? AND travel_date < ?", trainID, start, end).
		Delete(&model.Inventory{}).Error; err != nil {
		return err
	}

	inventories := make([]model.Inventory, 0, len(stops)*len(stops)*len(seeds)*7)
	for fromIndex := 0; fromIndex < len(stops)-1; fromIndex++ {
		for toIndex := fromIndex + 1; toIndex < len(stops); toIndex++ {
			fromStation := stations[stops[fromIndex].StationCode]
			toStation := stations[stops[toIndex].StationCode]
			mileage := stops[toIndex].Mileage - stops[fromIndex].Mileage
			if mileage <= 0 {
				continue
			}
			for day := 0; day < 7; day++ {
				travelDate := start.AddDate(0, 0, day)
				for _, seed := range seeds {
					totalCount := seed.TotalCount - day
					if totalCount < 1 {
						totalCount = 1
					}
					availableCount := demoAvailableCount(train.TrainNo, seed.SeatClassCode, seed.AvailableCount, totalCount, day)
					inventories = append(inventories, model.Inventory{
						TrainID:        trainID,
						TravelDate:     travelDate,
						FromStationID:  fromStation.ID,
						ToStationID:    toStation.ID,
						SeatClassCode:  seed.SeatClassCode,
						PriceCents:     seedFareCents(mileage, train.TrainType, seed.SeatClassCode),
						TotalCount:     totalCount,
						AvailableCount: availableCount,
						LockedCount:    0,
						SoldCount:      totalCount - availableCount,
						Status:         model.InventoryStatusActive,
					})
				}
			}
		}
	}
	if len(inventories) == 0 {
		return nil
	}
	return tx.CreateInBatches(inventories, 1000).Error
}

func demoAvailableCount(trainNo, seatClassCode string, baseAvailable, totalCount, day int) int {
	if totalCount <= 0 {
		return 0
	}

	trainWeight := 0
	for _, char := range trainNo {
		if char >= '0' && char <= '9' {
			trainWeight += int(char - '0')
		}
	}
	seatWeight := len(seatClassCode) % 5
	variation := (day+1)*(trainWeight%7+3) + seatWeight*2
	availableCount := baseAvailable - variation

	if day%3 == 0 && seatWeight == 0 {
		availableCount = 0
	}
	if day%4 == 1 && (seatClassCode == "BUSINESS" || seatClassCode == "DELUXE_SOFT_SLEEPER") {
		availableCount = 0
	}
	if day%5 == 2 {
		availableCount = minInt(3+trainWeight%4, totalCount)
	}
	if availableCount < 0 {
		availableCount = 0
	}
	if availableCount > totalCount {
		availableCount = totalCount
	}
	return availableCount
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func seedFareCents(mileage int, trainType string, seatClassCode string) int64 {
	coefficients := map[string]map[string]float64{
		"G": {"BUSINESS": 10.0, "FIRST": 5.8, "SECOND": 3.5},
		"C": {"BUSINESS": 10.0, "FIRST": 5.8, "SECOND": 3.5},
		"D": {"FIRST_SLEEPER": 5.0, "SECOND_SLEEPER": 3.8, "SECOND": 3.0},
		"Z": {"DELUXE_SOFT_SLEEPER": 4.0, "SOFT_SLEEPER": 3.0, "HARD_SLEEPER": 2.0, "HARD_SEAT": 1.0, "NO_SEAT": 1.0},
		"T": {"DELUXE_SOFT_SLEEPER": 4.0, "SOFT_SLEEPER": 3.0, "HARD_SLEEPER": 2.0, "HARD_SEAT": 1.0, "NO_SEAT": 1.0},
		"K": {"DELUXE_SOFT_SLEEPER": 4.0, "SOFT_SLEEPER": 3.0, "HARD_SLEEPER": 2.0, "HARD_SEAT": 1.0, "NO_SEAT": 1.0},
	}
	coefficient := coefficients[trainType][seatClassCode]
	if coefficient == 0 {
		coefficient = 1
	}
	return int64(math.Round(float64(mileage)*13*coefficient/100) * 100)
}
