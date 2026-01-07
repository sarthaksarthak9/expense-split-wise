package handlers

import (
	"expense-split-wise/internal/models"
	"expense-split-wise/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpenseHandler struct {
	expenseService *services.ExpenseService
	balanceService *services.BalanceService
}

func NewExpenseHandler(expenseService *services.ExpenseService, balanceService *services.BalanceService) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
		balanceService: balanceService,
	}
}

// CreateExpense handles POST /groups/:id/expenses
// Request: {"description": "Dinner", "amount": 1500, "paidBy": "Alice", "splitBetween": ["Alice", "Bob", "Charlie"]}
// Response: {"id": "...", "groupId": "...", "description": "Dinner", ...}
func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	groupID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var req struct {
		Description  string   `json:"description" binding:"required"`
		Amount       float64  `json:"amount" binding:"required,gt=0"`
		PaidBy       string   `json:"paidBy" binding:"required"`
		SplitBetween []string `json:"splitBetween" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense := &models.Expense{
		GroupID:      groupID,
		Description:  req.Description,
		Amount:       req.Amount,
		PaidBy:       req.PaidBy,
		SplitBetween: req.SplitBetween,
	}

	if err := h.expenseService.CreateExpense(c.Request.Context(), expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}

	c.JSON(http.StatusCreated, expense)
}

// GetExpenses handles GET /groups/:id/expenses
// Response: [{"id": "...", "description": "Dinner", ...}, ...]
func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	groupID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	expenses, err := h.expenseService.GetExpensesByGroup(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

// GetBalances handles GET /groups/:id/balances
// Response: {"Alice": 500, "Bob": -250, "Charlie": -250}
func (h *ExpenseHandler) GetBalances(c *gin.Context) {
	groupID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	balances, err := h.balanceService.GetBalances(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch balances"})
		return
	}

	c.JSON(http.StatusOK, balances)
}
