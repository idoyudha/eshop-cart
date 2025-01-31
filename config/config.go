package config

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"log"`
		MySQL
		Redis
		AuthService
		Kafka
		OrderService
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	MySQL struct {
		URL                   string `env-required:"true" env:"MYSQL_URL"`
		ConnectionMaxLifetime int    `env-required:"true" env:"MYSQL_CONN_MAX_LIFETIME"`
		MaxOpenConnection     int    `env-required:"true" env:"MYSQL_MAX_OPEN_CONNECTION"`
		MaxIdleConnection     int    `env-required:"true" env:"MYSQL_MAX_IDLE_CONNECTION"`
	}

	Redis struct {
		// RedisURL           string `env-required:"true" env:"REDIS_URL"`
		RedisMaster        string `env-required:"true" env:"REDIS_MASTER"`
		RedisSentinelAddrs string `env-required:"true" env:"REDIS_SENTINEL_ADDRS"`
		RedisPassword      string `env-required:"true" env:"REDIS_PASSWORD"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	AuthService struct {
		BaseURL string `env-required:"true" env:"AUTH_SERVICE"`
	}

	Kafka struct {
		Broker string `env-required:"true" env:"KAFKA_BROKER"`
	}

	OrderService struct {
		BaseURL string `env-required:"true" env:"ORDER_SERVICE"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
