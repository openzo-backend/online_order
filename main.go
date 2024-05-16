package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/online_order/config"
	handlers "github.com/tanush-128/openzo_backend/online_order/internal/api"
	"github.com/tanush-128/openzo_backend/online_order/internal/middlewares"
	"github.com/tanush-128/openzo_backend/online_order/internal/pb"
	"github.com/tanush-128/openzo_backend/online_order/internal/repository"
	"github.com/tanush-128/openzo_backend/online_order/internal/service"
	"google.golang.org/grpc"
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

	// // Initialize gRPC server
	// grpcServer := grpc.NewServer()
	// OnlineOrderpb.RegisterOnlineOrderServiceServer(grpcServer, service.NewGrpcOnlineOrderService(OnlineOrderRepository, OnlineOrderService))
	// reflection.Register(grpcServer) // Optional for server reflection

	//Initialize user gRPC client
	conn, err := grpc.Dial(cfg.UserGrpc, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)
	UserClient = c

	// initialize notification gRPC client
	notificationConn, err := grpc.Dial(cfg.NotificationGrpc, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer notificationConn.Close()
	notificationClient := pb.NewNotificationServiceClient(notificationConn)

	storeConn, err := grpc.Dial(cfg.StoreGrpc, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer storeConn.Close()
	storeClient := pb.NewStoreServiceClient(storeConn)

	online_orderRepository := repository.NewOnlineOrderRepository(db)
	OnlineOrderService := service.NewOnlineOrderService(online_orderRepository, notificationClient, storeClient)
	// Initialize HTTP server with Gin
	router := gin.Default()
	handler := handlers.NewHandler(&OnlineOrderService)

	router.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// router.Use(middlewares.JwtMiddleware(c))
	router.POST("/", handler.CreateOnlineOrder)
	router.GET("/:id", handler.GetOnlineOrderByID)
	router.GET("/changeStatus/:id", handler.ChangeOrderStatus)
	router.GET("/store/:store_id", handler.GetOnlineOrdersByStoreID)

	router.Use(middlewares.NewMiddleware(c).JwtMiddleware)
	router.PUT("/", handler.UpdateOnlineOrder)

	// router.Use(middlewares.JwtMiddleware)

	router.Run(fmt.Sprintf(":%s", cfg.HTTPPort))

}
