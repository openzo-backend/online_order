package repository

import (
	"github.com/google/uuid"
	"github.com/tanush-128/openzo_backend/online_order/internal/models"

	"gorm.io/gorm"
)

type OnlineOrderRepository interface {
	CreateOnlineOrder(OnlineOrder models.OnlineOrder) (models.OnlineOrder, error)
	GetOnlineOrderByID(id string) (models.OnlineOrder, error)
	GetOnlineOrdersByStoreID(store_id string) ([]models.OnlineOrder, error)

	UpdateOnlineOrder(OnlineOrder models.OnlineOrder) (models.OnlineOrder, error)
	// Add more methods for other OnlineOrder operations (GetOnlineOrderByEmail, UpdateOnlineOrder, etc.)

}

type online_orderRepository struct {
	db *gorm.DB
}

func NewOnlineOrderRepository(db *gorm.DB) OnlineOrderRepository {

	return &online_orderRepository{db: db}
}

func (r *online_orderRepository) CreateOnlineOrder(OnlineOrder models.OnlineOrder) (models.OnlineOrder, error) {
	OnlineOrder.ID = uuid.New().String()
	tx := r.db.Create(&OnlineOrder)

	if tx.Error != nil {
		return models.OnlineOrder{}, tx.Error
	}

	return OnlineOrder, nil
}

func (r *online_orderRepository) GetOnlineOrderByID(id string) (models.OnlineOrder, error) {
	var OnlineOrder models.OnlineOrder
	tx := r.db.Where("id = ?", id).Preload("Customer").Preload("OrderItems").First(&OnlineOrder)
	if tx.Error != nil {
		return models.OnlineOrder{}, tx.Error
	}

	return OnlineOrder, nil
}

func (r *online_orderRepository) GetOnlineOrdersByStoreID(store_id string) ([]models.OnlineOrder, error) {
	var OnlineOrders []models.OnlineOrder
	tx := r.db.Where("store_id = ?", store_id).Preload("OrderItems").Find(&OnlineOrders)
	if tx.Error != nil {
		return []models.OnlineOrder{}, tx.Error

	}

	return OnlineOrders, nil
}

func (r *online_orderRepository) UpdateOnlineOrder(OnlineOrder models.OnlineOrder) (models.OnlineOrder, error) {
	tx := r.db.Save(&OnlineOrder)
	if tx.Error != nil {
		return models.OnlineOrder{}, tx.Error
	}

	return OnlineOrder, nil
}

// Implement other repository methods (GetOnlineOrderByID, GetOnlineOrderByEmail, UpdateOnlineOrder, etc.) with proper error handling
