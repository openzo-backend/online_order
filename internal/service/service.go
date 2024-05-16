package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/online_order/internal/models"
	"github.com/tanush-128/openzo_backend/online_order/internal/pb"
	"github.com/tanush-128/openzo_backend/online_order/internal/repository"
)

type OnlineOrderService interface {

	//CRUD
	CreateOnlineOrder(ctx *gin.Context, req models.OnlineOrder) (models.OnlineOrder, error)
	GetOnlineOrderByID(ctx *gin.Context, id string) (models.OnlineOrder, error)
	GetOnlineOrdersByStoreID(ctx *gin.Context, store_id string) ([]models.OnlineOrder, error)
	ChangeOrderStatus(ctx *gin.Context, id string, status string) (models.OnlineOrder, error)
	UpdateOnlineOrder(ctx *gin.Context, req models.OnlineOrder) (models.OnlineOrder, error)
}

type online_orderService struct {
	online_orderRepository repository.OnlineOrderRepository
	notificationService    pb.NotificationServiceClient
	storeService           pb.StoreServiceClient
}

func NewOnlineOrderService(online_orderRepository repository.OnlineOrderRepository,
	notificationService pb.NotificationServiceClient, storeService pb.StoreServiceClient,
) OnlineOrderService {
	return &online_orderService{online_orderRepository: online_orderRepository, notificationService: notificationService, storeService: storeService}
}

func (s *online_orderService) CreateOnlineOrder(ctx *gin.Context, req models.OnlineOrder) (models.OnlineOrder, error) {

	token, err := s.storeService.GetFCMToken(ctx, &pb.StoreId{Id: req.StoreID})
	if err != nil {
		return models.OnlineOrder{}, err
	}

	_, err = s.notificationService.SendNotification(ctx, &pb.Notification{
		Title: "New Order",
		Body:  "Your store has a new order!",
		Token: token.Token,
		// Token:     "eFqUJcuhRWKFivWqRzeui3:APA91bHL8rpzdfk2dctBqIEHsdmiu8JJvTZgiF644vo39PaJ0g-yoc_BCSrSlxk7fTwznHvU4CaEup3CKA6FxXp-ZdS8TCzgWXOHDS504UAJQ9W5-l62V9gtOLkjnMXPKtHXeh5ZUp42",
		ImageURL:  "",
		ActionURL: "",
	})
	if err != nil {
		return models.OnlineOrder{}, err
	}
	createdOnlineOrder, err := s.online_orderRepository.CreateOnlineOrder(req)
	if err != nil {
		return models.OnlineOrder{}, err // Propagate error
	}

	return createdOnlineOrder, nil
}

func (s *online_orderService) ChangeOrderStatus(ctx *gin.Context, id string, status string) (models.OnlineOrder, error) {

	token := ""

	if status == "accepted" {
		acceptedNotification.Token = token
		_, err := s.notificationService.SendNotification(ctx, &acceptedNotification)
		if err != nil {
			return models.OnlineOrder{}, err
		}
	} else if status == "rejected" {
		rejectedNotification.Token = token
		_, err := s.notificationService.SendNotification(ctx, &rejectedNotification)
		if err != nil {
			return models.OnlineOrder{}, err
		}

	} else if status == "out_for_delivery" {
		outForDeliveryNotification.Token = token
		_, err := s.notificationService.SendNotification(ctx, &outForDeliveryNotification)
		if err != nil {
			return models.OnlineOrder{}, err
		}

	} else if status == "cancelled" {
		cancelledNotification.Token = token
		_, err := s.notificationService.SendNotification(ctx, &cancelledNotification)
		if err != nil {
			return models.OnlineOrder{}, err
		}

	} else if status == "delivered" {
		deliveredNotification.Token = token
		_, err := s.notificationService.SendNotification(ctx, &deliveredNotification)
		if err != nil {
			return models.OnlineOrder{}, err
		}
	}

	changedOnlineOrder, err := s.online_orderRepository.ChangeOrderStatus(id, status)
	if err != nil {
		return models.OnlineOrder{}, err
	}

	return changedOnlineOrder, nil
}

func (s *online_orderService) UpdateOnlineOrder(ctx *gin.Context, req models.OnlineOrder) (models.OnlineOrder, error) {
	updatedOnlineOrder, err := s.online_orderRepository.UpdateOnlineOrder(req)
	if err != nil {
		return models.OnlineOrder{}, err
	}

	return updatedOnlineOrder, nil
}
