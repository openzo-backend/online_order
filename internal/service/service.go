package service

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/online_order/internal/models"
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
	kafkaProducer          *kafka.Producer
}

func NewOnlineOrderService(online_orderRepository repository.OnlineOrderRepository,
	kafkaProducer *kafka.Producer,
) OnlineOrderService {
	return &online_orderService{online_orderRepository: online_orderRepository, kafkaProducer: kafkaProducer}
}

func (s *online_orderService) CreateOnlineOrder(ctx *gin.Context, req models.OnlineOrder) (models.OnlineOrder, error) {

	topic := "onlineorder"
	orderMsg, err := json.Marshal(req)
	if err != nil {
		return models.OnlineOrder{}, err
	}

	s.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          orderMsg,
		Key:            []byte("new_order"),
	}, nil)
	s.kafkaProducer.Flush(15 * 1000)

	createdOnlineOrder, err := s.online_orderRepository.CreateOnlineOrder(req)
	if err != nil {
		return models.OnlineOrder{}, err // Propagate error
	}

	return createdOnlineOrder, nil
}

func (s *online_orderService) ChangeOrderStatus(ctx *gin.Context, id string, status string) (models.OnlineOrder, error) {

	OnlineOrder, err := s.online_orderRepository.GetOnlineOrderByID(id)
	if err != nil {
		return models.OnlineOrder{}, err
	}
	if status == "accepted" {
		OnlineOrder.OrderStatus = models.OrderAccepted
	} else if status == "rejected" {
		OnlineOrder.OrderStatus = models.OrderRejected
	} else if status == "out_for_delivery" {
		OnlineOrder.OrderStatus = models.OrderOutForDel
	} else if status == "cancelled" {
		OnlineOrder.OrderStatus = models.OrderCancelled
	} else if status == "delivered" {
		OnlineOrder.OrderStatus = models.OrderDelivered
	}

	topic := "order-status-updates"
	orderStatusMsg, err := json.Marshal(OnlineOrder)
	if err != nil {
		return models.OnlineOrder{}, err
	}

	s.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          orderStatusMsg,
		Key:            []byte("order_status"),
	}, nil)
	s.kafkaProducer.Flush(15 * 1000)

	changedOnlineOrder, err := s.online_orderRepository.UpdateOnlineOrder(OnlineOrder)
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
