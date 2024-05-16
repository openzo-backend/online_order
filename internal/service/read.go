package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/online_order/internal/models"
)

func (s *online_orderService) GetOnlineOrderByID(ctx *gin.Context, id string) (models.OnlineOrder, error) {
	online_order, err := s.online_orderRepository.GetOnlineOrderByID(id)
	if err != nil {
		return models.OnlineOrder{}, err
	}

	return online_order, nil
}

func (s *online_orderService) GetOnlineOrdersByStoreID(ctx *gin.Context, store_id string) ([]models.OnlineOrder, error) {
	online_orders, err := s.online_orderRepository.GetOnlineOrdersByStoreID(store_id)
	if err != nil {
		return []models.OnlineOrder{}, err
	}

	return online_orders, nil
}

