package order

import (
	"dz4/internal/models"
	"dz4/pkg/db"

	"gorm.io/gorm/clause"
)

type OrderRepository struct {
	Database *db.Db
}

func NewOrderRepository(db *db.Db) *OrderRepository {
	return &OrderRepository{
		Database: db,
	}
}

func (repo *OrderRepository) Create(order *models.Order) (*models.Order, error) {
	result := repo.Database.DB.Create(order)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}

func (repo *OrderRepository) GetById(id uint) (*models.Order, error) {
	var order models.Order
	result := repo.Database.DB.Preload("Products").First(&order, id)
	// Метод Preload("Products") указывает GORM:
	// «Когда забираешь заказы (orders), сразу подгрузи к каждому все связанные продукты из таблицы products.»
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

func (repo *OrderRepository) GetByUserId(userId uint) ([]models.Order, error) {
	var orders []models.Order
	result := repo.Database.DB.Preload("Products").Where("user_id = ?", userId).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

func (repo *OrderRepository) AddProductsToOrder(orderId uint, productIds []uint) error {
	var order models.Order
	if err := repo.Database.DB.First(&order, orderId).Error; err != nil {
		return err
	}

	var products []models.Product
	if err := repo.Database.DB.Where("id IN ?", productIds).Find(&products).Error; err != nil {
		return err
	}

	if err := repo.Database.DB.Model(&order).Association("Products").Append(products); err != nil {
		return err
	}

	return nil
}

func (repo *OrderRepository) Update(order *models.Order) (*models.Order, error) {
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(order)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}

func (repo *OrderRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(&models.Order{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}