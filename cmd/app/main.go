package app

import (
	"log"

	"github.com/idoyudha/eshop-cart/config"
	"github.com/idoyudha/eshop-cart/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Config error: ", err)
	}

	app.Run(cfg)
}
