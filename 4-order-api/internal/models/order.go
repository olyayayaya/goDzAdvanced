package models

import (

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserId   uint              `json:"user_id"`
	Date     datatypes.Date    `json:"date"`
	Products []Product `json:"products" gorm:"many2many:order_products;"`
}
