package services

import (
	"context"
	"encoding/json"
	"expense-split-wise/internal/database"
	"expense-split-wise/internal/models"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BalanceService struct {
	mongo *database.MongoClient
	redis *database.RedisClient
}

func NewBalanceService(mongo *database.MongoClient, redis *database.RedisClient) *BalanceService {
	return &BalanceService{
		mongo: mongo,
		redis: redis,
	}
}

// RecalculateBalances recalculates balances for a group
func (s *BalanceService) RecalculateBalances(ctx context.Context, groupID primitive.ObjectID) error {
	// Fetch all expenses for the group
	cursor, err := s.mongo.Collection("expenses").Find(ctx, bson.M{"groupId": groupID})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var expenses []models.Expense
	if err := cursor.All(ctx, &expenses); err != nil {
		return err
	}

	// Calculate balances
	balances := make(map[string]float64)

	for _, expense := range expenses {
		// Split amount equally among splitBetween members
		splitAmount := expense.Amount / float64(len(expense.SplitBetween))

		// Add to payer's balance (they are owed)
		balances[expense.PaidBy] += expense.Amount

		// Deduct from each member's balance (they owe)
		for _, member := range expense.SplitBetween {
			balances[member] -= splitAmount
		}
	}

	// Save to MongoDB
	balance := &models.Balance{
		GroupID:   groupID,
		Balances:  balances,
		UpdatedAt: time.Now(),
	}

	opts := options.Update().SetUpsert(true)
	_, err = s.mongo.Collection("balances").UpdateOne(
		ctx,
		bson.M{"groupId": groupID},
		bson.M{"$set": balance},
		opts,
	)
	if err != nil {
		return err
	}

	// Cache in Redis for fast retrieval
	cacheKey := fmt.Sprintf("balance:%s", groupID.Hex())
	data, _ := json.Marshal(balances)
	s.redis.Client.Set(ctx, cacheKey, data, 30*time.Minute)

	return nil
}

// GetBalances retrieves balances for a group (from cache or DB)
func (s *BalanceService) GetBalances(ctx context.Context, groupID primitive.ObjectID) (map[string]float64, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("balance:%s", groupID.Hex())
	cached, err := s.redis.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		var balances map[string]float64
		if json.Unmarshal([]byte(cached), &balances) == nil {
			return balances, nil
		}
	}

	// Fallback to database
	var balance models.Balance
	err = s.mongo.Collection("balances").FindOne(ctx, bson.M{"groupId": groupID}).Decode(&balance)
	if err != nil {
		return nil, err
	}

	// Update cache
	data, _ := json.Marshal(balance.Balances)
	s.redis.Client.Set(ctx, cacheKey, data, 30*time.Minute)

	return balance.Balances, nil
}
