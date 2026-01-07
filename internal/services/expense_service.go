package services

import (
	"context"
	"expense-split-wise/internal/models"
	"expense-split-wise/internal/queue"
	"expense-split-worker/internal/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpenseService struct {
	mongo    *database.MongoClient
	rabbitmq *queue.RabbitMQClient
	queue    string
}

func NewExpenseService(mongo *database.MongoClient, rabbitmq *queue.RabbitMQClient, queueName string) *ExpenseService {
	return &ExpenseService{
		mongo:    mongo,
		rabbitmq: rabbitmq,
		queue:    queueName,
	}
}

// CreateExpense creates a new expense and publishes to queue
func (s *ExpenseService) CreateExpense(ctx context.Context, expense *models.Expense) error {
	expense.CreatedAt = time.Now()

	// Save expense to MongoDB
	result, err := s.mongo.Collection("expenses").InsertOne(ctx, expense)
	if err != nil {
		return err
	}

	expense.ID = result.InsertedID.(primitive.ObjectID)

	// Publish message to RabbitMQ for async processing
	message := models.ExpenseMessage{
		GroupID:   expense.GroupID.Hex(),
		ExpenseID: expense.ID.Hex(),
		Amount:    expense.Amount,
	}

	return s.rabbitmq.PublishMessage(s.queue, message)
}

// GetExpensesByGroup retrieves all expenses for a group
func (s *ExpenseService) GetExpensesByGroup(ctx context.Context, groupID primitive.ObjectID) ([]models.Expense, error) {
	cursor, err := s.mongo.Collection("expenses").Find(ctx, bson.M{"groupId": groupID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var expenses []models.Expense
	if err := cursor.All(ctx, &expenses); err != nil {
		return nil, err
	}

	return expenses, nil
}
