package main

import (
	"github.com/tomohavvk/go-walker/configs"
	"github.com/tomohavvk/go-walker/db"
)

func main() {
	cfg := configs.LoadConfig()

	err := db.PerformMigration(cfg.DB)
	if err != nil {
		panic(err)
	}
}
