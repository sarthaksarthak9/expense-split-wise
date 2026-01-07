package main

import (
	"expense-split-wise/internal/config"
	"expense-split-wise/internal/database"
	"expense-split-wise/internal/handlers"
	"expense-split-wise/internal/queue"
	"expense-split-wise/internal/services"
	"log"

	"github.com/gin-gonic/gin"
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

	// Declare the expense queue
	if err := rabbitmq.DeclareQueue(cfg.ExpenseQueue); err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Initialize services
	groupService := services.NewGroupService(mongoDB)
	expenseService := services.NewExpenseService(mongoDB, rabbitmq, cfg.ExpenseQueue)
	balanceService := services.NewBalanceService(mongoDB, redisClient)

	// Initialize handlers
	groupHandler := handlers.NewGroupHandler(groupService)
	expenseHandler := handlers.NewExpenseHandler(expenseService, balanceService)

	// Setup Gin router
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Group routes
		api.POST("/groups", groupHandler.CreateGroup)
		api.POST("/groups/:id/users", groupHandler.AddUsersToGroup)

		// Expense routes
		api.POST("/groups/:id/expenses", expenseHandler.CreateExpense)
		api.GET("/groups/:id/expenses", expenseHandler.GetExpenses)

		// Balance routes
		api.GET("/groups/:id/balances", expenseHandler.GetBalances)
	}

	// Start server
	log.Printf("🚀 API Server starting on port %s", cfg.APIPort)
	if err := router.Run(":" + cfg.APIPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
