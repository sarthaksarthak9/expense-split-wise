package handlers

import (
	"expense-split-worker/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupHandler struct {
	groupService *services.GroupService
}

func NewGroupHandler(groupService *services.GroupService) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

// CreateGroup handles POST /groups
// Request: {"name": "Trip to Goa", "members": ["Alice", "Bob", "Charlie"]}
// Response: {"id": "...", "name": "Trip to Goa", "members": [...]}
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req struct {
		Name    string   `json:"name" binding:"required"`
		Members []string `json:"members" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := h.groupService.CreateGroup(c.Request.Context(), req.Name, req.Members)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// AddUsersToGroup handles POST /groups/:id/users
// Request: {"members": ["David", "Eve"]}
// Response: {"message": "Users added successfully"}
func (h *GroupHandler) AddUsersToGroup(c *gin.Context) {
	groupID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var req struct {
		Members []string `json:"members" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.groupService.AddMembersToGroup(c.Request.Context(), groupID, req.Members); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Users added successfully"})
}
