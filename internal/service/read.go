package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/online_order/internal/models"
)

func (s *online_orderService) GetOnlineOrderByID(ctx *gin.Context, id string) (models.OnlineOrder, error) {

	return s.online_orderRepository.GetOnlineOrderByID(id)
}
func (s *online_orderService) GetOnlineOrdersByStoreID(ctx *gin.Context, store_id string) ([]models.OnlineOrder, error) {

	return s.online_orderRepository.GetOnlineOrdersByStoreID(store_id)

}

func (s *online_orderService) GetOnlineOrdersByUserDataId(ctx *gin.Context, user_data_id string) ([]models.OnlineOrder, error) {

	return s.online_orderRepository.GetOnlineOrdersByUserDataId(user_data_id)
}
