package main

import (
	"log"

	"github.com/idoyudha/eshop-cart/config"
	"github.com/idoyudha/eshop-cart/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
