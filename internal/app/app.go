package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-cart/config"
	v1 "github.com/idoyudha/eshop-cart/internal/controller/http/v1"
	"github.com/idoyudha/eshop-cart/internal/usecase"
	"github.com/idoyudha/eshop-cart/internal/usecase/repo"
	"github.com/idoyudha/eshop-cart/pkg/httpserver"
	"github.com/idoyudha/eshop-cart/pkg/logger"
	"github.com/idoyudha/eshop-cart/pkg/mysql"
	"github.com/idoyudha/eshop-cart/pkg/redis"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	mySQL, err := mysql.NewMySQL(cfg.MySQL)
	if err != nil {
		l.Fatal("app - Run - dynamodb.NewDynamoDB: ", err)
	}

	redisClient := redis.NewRedis(cfg.Redis)

	cartUseCase := usecase.NewCartUseCase(
		repo.NewCartRedisRepo(redisClient),
		repo.NewCartMySQLRepo(mySQL),
	)

	handler := gin.Default()
	v1.NewRouter(handler, cartUseCase, l)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error("app - Run - httpServer.Notify: ", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Info("app - Run - httpServer.Shutdown: %s", err)
	}
}
