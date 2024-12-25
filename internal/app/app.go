package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-cart/config"
	v1Http "github.com/idoyudha/eshop-cart/internal/controller/http/v1"
	v1Kafka "github.com/idoyudha/eshop-cart/internal/controller/kafka/v1"
	"github.com/idoyudha/eshop-cart/internal/usecase"
	"github.com/idoyudha/eshop-cart/internal/usecase/repo"
	"github.com/idoyudha/eshop-cart/pkg/httpserver"
	"github.com/idoyudha/eshop-cart/pkg/kafka"
	"github.com/idoyudha/eshop-cart/pkg/logger"
	"github.com/idoyudha/eshop-cart/pkg/mysql"
	"github.com/idoyudha/eshop-cart/pkg/redis"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	kafkaConsumer, err := kafka.NewKafkaConsumer(cfg.Kafka.Broker)
	if err != nil {
		l.Fatal("app - Run - kafka.NewKafkaConsumer: ", err)
	}
	defer kafkaConsumer.Close()

	mySQL, err := mysql.NewMySQL(cfg.MySQL)
	if err != nil {
		l.Fatal("app - Run - mysql.NewMySQL: ", err)
	}

	redisClient, err := redis.NewRedis(cfg.Redis)
	if err != nil {
		l.Fatal("app - Run - redis.NewRedis: ", err)
	}

	cartUseCase := usecase.NewCartUseCase(
		repo.NewCartRedisRepo(redisClient),
		repo.NewCartMySQLRepo(mySQL),
		cfg.OrderService,
	)

	// HTTP Server
	handler := gin.Default()
	v1Http.NewRouter(handler, cartUseCase, l, cfg.AuthService)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Kafka Consumer
	kafkaErrChan := make(chan error, 1)
	go func() {
		if err := v1Kafka.KafkaNewRouter(cartUseCase, l, kafkaConsumer); err != nil {
			kafkaErrChan <- err
		}
	}()

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
