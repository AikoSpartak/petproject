package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Sslmode  string
}

func LoadConfig() (*Config, error) {

	//загруждаем файл
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading file")
	}

	//получаем значение переменной из окружения
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	//делаем валидацию
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		return nil, errors.New("Error validations")
	}

	cfg := Config{
		DBConfig: DBConfig{
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			Name:     dbname,
			Sslmode:  sslmode,
		},
	}

	return &cfg, nil

}

// создали строку подключения
func (c *Config) ConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBConfig.Host, c.DBConfig.Port, c.DBConfig.User, c.DBConfig.Password, c.DBConfig.Name, c.DBConfig.Sslmode)
}
