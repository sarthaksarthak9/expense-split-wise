package main

import (
	"expense-split-wise/internal/config"
	"expense-split-wise/internal/database"
	"expense-split-wise/internal/queue"
	"expense-split-wise/internal/services"
	"expense-split-wise/internal/worker"
	"log"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize MongoDB
	mongoDB, err := database.NewMongoClient(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	// Initialize Redis
	redisClient, err := database.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize RabbitMQ
	rabbitmq, err := queue.NewRabbitMQClient(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	// Initialize services
	balanceService := services.NewBalanceService(mongoDB, redisClient)

	// Initialize and start worker
	expenseWorker := worker.NewExpenseWorker(rabbitmq, balanceService, cfg.ExpenseQueue)

	log.Println("🚀 Worker starting...")
	if err := expenseWorker.Start(); err != nil {
		log.Fatalf("Worker failed: %v", err)
	}
}
