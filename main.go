package main

import (
	"log"

	"github.com/e-commerce/common/config"
	"github.com/e-commerce/common/constant"
	"github.com/e-commerce/common/database"
	"github.com/e-commerce/model"

	"github.com/google/gops/agent"
)

func main() {
	log.Println("main started.......")

	opts := agent.Options{
		ShutdownCleanup: true,
	}
	if err := agent.Listen(opts); err != nil {
		log.Fatal(err)
	}

	cfg := config.GetConfig()

	initDatabase(cfg.Database)

}

func initDatabase(dbConfigs map[string]*config.DatabaseConfig) {
	log.Println("Database Initialized....... ")
	database.InitDatabase(dbConfigs)

	ecommercedb := database.DBConnMap[constant.DB_NAME]
	if ecommercedb == nil {
		log.Fatalf("Unable to initialize OMS DB")
	} else {
		model.Init(ecommercedb)
	}

}
