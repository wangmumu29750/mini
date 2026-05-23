package main

import (
	"log"

	"mini-12306/backend/internal/config"
	"mini-12306/backend/internal/database"
	"mini-12306/backend/internal/router"
)

func main() {
	cfg := config.Load()
//cl到此一游
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate database: %v", err)
	}
	if err := database.SeedDemoData(db); err != nil {
		log.Fatalf("seed demo data: %v", err)
	}

	r := router.New(cfg, db)

	log.Printf("mini-12306 backend listening on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
