package services

import (
	"context"
	"expense-split-wise/internal/models"
	"expense-split-worker/internal/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupService struct {
	mongo *database.MongoClient
}

func NewGroupService(mongo *database.MongoClient) *GroupService {
	return &GroupService{mongo: mongo}
}

// CreateGroup creates a new group
func (s *GroupService) CreateGroup(ctx context.Context, name string, members []string) (*models.Group, error) {
	group := &models.Group{
		Name:      name,
		Members:   members,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := s.mongo.Collection("groups").InsertOne(ctx, group)
	if err != nil {
		return nil, err
	}

	group.ID = result.InsertedID.(primitive.ObjectID)
	return group, nil
}

// GetGroup retrieves a group by ID
func (s *GroupService) GetGroup(ctx context.Context, id primitive.ObjectID) (*models.Group, error) {
	var group models.Group
	err := s.mongo.Collection("groups").FindOne(ctx, bson.M{"_id": id}).Decode(&group)
	return &group, err
}

// AddMembersToGroup adds users to an existing group
func (s *GroupService) AddMembersToGroup(ctx context.Context, groupID primitive.ObjectID, members []string) error {
	update := bson.M{
		"$addToSet": bson.M{"members": bson.M{"$each": members}},
		"$set":      bson.M{"updatedAt": time.Now()},
	}
	_, err := s.mongo.Collection("groups").UpdateOne(ctx, bson.M{"_id": groupID}, update)
	return err
}
