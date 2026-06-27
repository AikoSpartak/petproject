package db

import (
	"arch/internal/config"
	"arch/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func FirstPostgresLoad(cfg *config.Config) (*gorm.DB, error) {

	//передаем строку подключения
	connStr := cfg.ConnString()

	// подключаем к дб
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	//достали sql.DB через gormDB, так как у grom нет прямого Ping.
	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	//проверяем на доступность базу
	err = sqlDb.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Подключение к БД успешно")

	//смотрим на структуру User и создаем таблицу users с такими же полями.
	//если таблица уже создана, то проверяем, нет ли новых полей. Если есть - добавляем
	db.AutoMigrate(&models.User{})

	return db, nil
}
