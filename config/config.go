package config

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		MySQL
		Redis
		Log `yaml:"log"`
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
		RedisURL      string `env-required:"true" env:"REDIS_URL"`
		RedisPassword string `env-required:"true" env:"REDIS_PASSWORD"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
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
