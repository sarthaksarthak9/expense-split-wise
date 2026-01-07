package worker

import (
	"context"
	"encoding/json"
	"expense-split-wise/internal/models"
	"expense-split-wise/internal/queue"
	"expense-split-wise/internal/services"
	"log"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpenseWorker struct {
	rabbitmq       *queue.RabbitMQClient
	balanceService *services.BalanceService
	queueName      string
}

func NewExpenseWorker(rabbitmq *queue.RabbitMQClient, balanceService *services.BalanceService, queueName string) *ExpenseWorker {
	return &ExpenseWorker{
		rabbitmq:       rabbitmq,
		balanceService: balanceService,
		queueName:      queueName,
	}
}

// Start begins consuming messages from the queue
func (w *ExpenseWorker) Start() error {
	// Declare the queue
	if err := w.rabbitmq.DeclareQueue(w.queueName); err != nil {
		return err
	}

	// Start consuming messages
	messages, err := w.rabbitmq.ConsumeMessages(w.queueName)
	if err != nil {
		return err
	}

	log.Printf("🔄 Worker started, listening on queue: %s", w.queueName)

	// Process messages continuously
	forever := make(chan bool)

	go func() {
		for msg := range messages {
			w.processMessage(msg)
		}
	}()

	<-forever
	return nil
}

// processMessage handles individual expense messages
func (w *ExpenseWorker) processMessage(msg amqp.Delivery) {
	var expenseMsg models.ExpenseMessage
	if err := json.Unmarshal(msg.Body, &expenseMsg); err != nil {
		log.Printf("❌ Failed to unmarshal message: %v", err)
		msg.Nack(false, false) // Reject and don't requeue
		return
	}

	log.Printf("📨 Processing expense: %s for group: %s", expenseMsg.ExpenseID, expenseMsg.GroupID)

	// Convert groupID string to ObjectID
	groupID, err := primitive.ObjectIDFromHex(expenseMsg.GroupID)
	if err != nil {
		log.Printf("❌ Invalid group ID: %v", err)
		msg.Nack(false, false)
		return
	}

	// Recalculate balances for the group
	ctx := context.Background()
	if err := w.balanceService.RecalculateBalances(ctx, groupID); err != nil {
		log.Printf("❌ Failed to recalculate balances: %v", err)
		msg.Nack(false, true) // Reject and requeue for retry
		return
	}

	log.Printf("✅ Successfully recalculated balances for group: %s", expenseMsg.GroupID)
	msg.Ack(false) // Acknowledge successful processing
}
