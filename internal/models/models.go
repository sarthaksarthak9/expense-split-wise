package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Group represents a group of users who share expenses
type Group struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Members   []string           `json:"members" bson:"members"` // User IDs or names
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// Expense represents a shared expense in a group
type Expense struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	GroupID      primitive.ObjectID `json:"groupId" bson:"groupId"`
	Description  string             `json:"description" bson:"description"`
	Amount       float64            `json:"amount" bson:"amount"`
	PaidBy       string             `json:"paidBy" bson:"paidBy"`             // User who paid
	SplitBetween []string           `json:"splitBetween" bson:"splitBetween"` // Users to split between
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
}

// Balance represents the balance sheet for a group
type Balance struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	GroupID   primitive.ObjectID `json:"groupId" bson:"groupId"`
	Balances  map[string]float64 `json:"balances" bson:"balances"` // user -> amount (positive = owed, negative = owes)
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// ExpenseMessage represents the message sent to RabbitMQ
type ExpenseMessage struct {
	GroupID   string  `json:"groupId"`
	ExpenseID string  `json:"expenseId"`
	Amount    float64 `json:"amount"`
}
