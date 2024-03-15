package config

import "fmt"

type ApiConfig struct {
	AppConfig AppConfig
	DbConfig  DbConfig
	Host      string `env:"API_HOST, required"`
	Port      string `env:"API_PORT, required"`
}

type AppConfig struct {
	Name     string `env:"APP_NAME, required"`
	Env      string `env:"APP_ENV, required"`
	LogLevel string `env:"APP_LOG_LEVEL, default=info"`
}

type DbConfig struct {
	Host    string `env:"DB_HOST, required"`
	Port    string `env:"DB_PORT, required"`
	User    string `env:"DB_USER, required"`
	Pass    string `env:"DB_PASS, required"`
	DbName  string `env:"DB_NAME, required"`
	SslMode string `env:"DB_SSLMODE, required"`
}

func (ac ApiConfig) RetrieveApiAddress() string {
	return fmt.Sprintf("%v:%v", ac.Host, ac.Port)
}

func (dc DbConfig) RetrieveDBConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dc.Host, dc.Port, dc.User, dc.Pass, dc.DbName,
	)
}
