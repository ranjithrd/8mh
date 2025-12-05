package handlers

import (
	"backend/src/db"
	"backend/src/repos"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

var userRepoHandler = repos.User{}

type UserListItem struct {
	ID          uint   `json:"user_id" example:"1"`
	Name        string `json:"name" example:"John Doe"`
	PhoneNumber string `json:"phone_number" example:"+1234567890"`
}

type UserListResponse struct {
	Users []UserListItem `json:"users"`
}

type UserDetailResponse struct {
	ID             uint   `json:"user_id" example:"1"`
	Name           string `json:"name" example:"John Doe"`
	PhoneNumber    string `json:"phone_number" example:"+1234567890"`
	SavingsBalance int    `json:"savings_balance" example:"50000"`
	SharesBalance  int    `json:"shares_balance" example:"25000"`
}

// ListUsers godoc
// @Summary List/Search Users
// @Description Get list of all users, optionally filter by name or phone number
// @Tags users
// @Produce json
// @Param search query string false "Search by name or phone number"
// @Success 200 {object} UserListResponse
// @Failure 500 {object} ErrorResponse
// @Security SessionAuth
// @Router /api/v1/users [get]
func ListUsers(c echo.Context) error {
	search := c.QueryParam("search")

	var users []db.User
	query := db.DB

	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR phone_number LIKE ?", searchTerm, searchTerm)
	}

	if err := query.Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
	}

	userList := make([]UserListItem, len(users))
	for i, user := range users {
		userList[i] = UserListItem{
			ID:          user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		}
	}

	return c.JSON(http.StatusOK, UserListResponse{Users: userList})
}

// GetUserByID godoc
// @Summary Get User Details
// @Description Get detailed information about a specific user including balances
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserDetailResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security SessionAuth
// @Router /api/v1/users/{id} [get]
func GetUserByID(c echo.Context) error {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var user db.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// Get savings balance
	var savingsBalance int
	db.DB.Model(&db.Deposit{}).
		Where("user_id = ? AND type = ?", userID, "savings").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&savingsBalance)

	// Get shares balance
	var sharesBalance int
	db.DB.Model(&db.Deposit{}).
		Where("user_id = ? AND type = ?", userID, "shares").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&sharesBalance)

	response := UserDetailResponse{
		ID:             user.ID,
		Name:           user.Name,
		PhoneNumber:    user.PhoneNumber,
		SavingsBalance: savingsBalance,
		SharesBalance:  sharesBalance,
	}

	return c.JSON(http.StatusOK, response)
}
