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
	GetOnlineOrdersByUserDataId(user_data_id string) ([]models.OnlineOrder, error)
	UpdateOnlineOrder(OnlineOrder models.OnlineOrder) (models.OnlineOrder, error)
	DeleteOnlineOrder(id string) error
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
	tx := r.db.Where("store_id = ?", store_id).Preload("OrderItems").Preload("Customer").Find(&OnlineOrders)
	if tx.Error != nil {
		return []models.OnlineOrder{}, tx.Error

	}

	return OnlineOrders, nil
}

//get all orders by where orders.customer.user_data_id = user_data_id

func (r *online_orderRepository) GetOnlineOrdersByUserDataId(user_data_id string) ([]models.OnlineOrder, error) {
	var OnlineOrders []models.OnlineOrder

	tx := r.db.Joins("Customer").Where("user_data_id = ?", user_data_id).Preload("OrderItems").Preload("Customer").Find(&OnlineOrders)
	if tx.Error != nil {
		return []models.OnlineOrder{}, tx.Error

	}

	return OnlineOrders, nil
}

func (r *online_orderRepository) UpdateOnlineOrder(OnlineOrder models.OnlineOrder) (models.OnlineOrder, error) {
	// update order items first
	tx := r.db.Where("online_order_id = ?", OnlineOrder.ID).Delete(&models.OnlineOrderItem{})
	if tx.Error != nil {
		return models.OnlineOrder{}, tx.Error
	}

	// update online customer
	tx = r.db.Where("online_order_id = ?", OnlineOrder.ID).Delete(&models.OnlineCustomer{})
	if tx.Error != nil {
		return models.OnlineOrder{}, tx.Error
	}

	// update online order

	tx = r.db.Save(&OnlineOrder)
	if tx.Error != nil {
		return models.OnlineOrder{}, tx.Error
	}

	return OnlineOrder, nil
}

func (r *online_orderRepository) DeleteOnlineOrder(id string) error {

	// delete order items first
	tx := r.db.Where("online_order_id = ?", id).Delete(&models.OnlineOrderItem{})
	if tx.Error != nil {
		return tx.Error
	}

	// delete online customer

	tx = r.db.Where("online_order_id = ?", id).Delete(&models.OnlineCustomer{})
	if tx.Error != nil {
		return tx.Error

	}

	// delete online order

	tx = r.db.Where("id = ?", id).Delete(&models.OnlineOrder{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// Implement other repository methods (GetOnlineOrderByID, GetOnlineOrderByEmail, UpdateOnlineOrder, etc.) with proper error handling
