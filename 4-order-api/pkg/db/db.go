package db

import (
	"dz4/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}
//структура, в которой хранится наша бд

func NewDb(conf *configs.Config) *Db {
	db, err := gorm.Open(postgres.Open(conf.Db.Dsn), &gorm.Config{}) // открываем соединение, чтобы подключиться в бд
	if err != nil {
		panic(err)
	}
	return &Db{db}
}
