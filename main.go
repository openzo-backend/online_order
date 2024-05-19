package main

import (
	"log"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/online_order/config"
	handlers "github.com/tanush-128/openzo_backend/online_order/internal/api"
	"github.com/tanush-128/openzo_backend/online_order/internal/middlewares"
	"github.com/tanush-128/openzo_backend/online_order/internal/pb"
	"github.com/tanush-128/openzo_backend/online_order/internal/repository"
	"github.com/tanush-128/openzo_backend/online_order/internal/service"
	"google.golang.org/grpc"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var UserClient pb.UserServiceClient

type User2 struct {
}

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %w", err))
	}

	db, err := connectToDB(cfg) // Implement database connection logic
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to database: %w", err))
	}

	// creates a new producer instance
	conf := ReadConfig()
	p, _ := kafka.NewProducer(&conf)

	// go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	//Initialize user gRPC client
	conn, err := grpc.Dial(cfg.UserGrpc, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)
	UserClient = c

	online_orderRepository := repository.NewOnlineOrderRepository(db)
	OnlineOrderService := service.NewOnlineOrderService(online_orderRepository, p)
	// Initialize HTTP server with Gin
	router := gin.Default()
	handler := handlers.NewHandler(&OnlineOrderService)

	router.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/", handler.CreateOnlineOrder)
	router.GET("/:id", handler.GetOnlineOrderByID)
	router.GET("/changeStatus/:id", handler.ChangeOrderStatus)
	router.GET("/store/:store_id", handler.GetOnlineOrdersByStoreID)
	router.Use(middlewares.NewMiddleware(c).JwtMiddleware)
	router.PUT("/", handler.UpdateOnlineOrder)

	router.Run(fmt.Sprintf(":%s", cfg.HTTPPort))

}
