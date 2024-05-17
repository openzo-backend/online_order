package main

import (
	"log"

	"bufio"
	"fmt"
	"os"
	"strings"

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

func ReadConfig() kafka.ConfigMap {
	// reads the client configuration from client.properties
	// and returns it as a key-value map
	m := make(map[string]kafka.ConfigValue)

	file, err := os.Open("client.properties")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %s", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			kv := strings.Split(line, "=")
			parameter := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			m[parameter] = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to read file: %s", err)
		os.Exit(1)
	}

	return m
}

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
	topic := "onlineorder"

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

	// produces a sample message to the user-created topic
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte("key"),
		Value:          []byte("value"),
	}, nil)

	// send any outstanding or buffered messages to the Kafka broker and close the connection
	p.Flush(15 * 1000)
	p.Close()
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
