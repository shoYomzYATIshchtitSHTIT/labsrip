package handler

import (
	"Backend-RIP/internal/app/ds"
	"Backend-RIP/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	repo *repository.Repository
}

func NewUserHandler(repo *repository.Repository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

type RegisterRequest struct {
	Login       string `json:"login" binding:"required"`
	Password    string `json:"password" binding:"required"`
	IsModerator bool   `json:"is_moderator"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Login    *string `json:"login"`
	Password *string `json:"password"`
}

// POST регистрация
func (h *UserHandler) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	user := &ds.Users{
		Login:       req.Login,
		Password:    req.Password,
		IsModerator: req.IsModerator,
	}

	if err := h.repo.User.RegisterUser(user); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": user.User_ID,
	})
}

// GET полей пользователя после аутентификации (для личного кабинета)
func (h *UserHandler) GetProfile(ctx *gin.Context) {
	// В реальном приложении userID брался бы из контекста аутентификации
	// Для примера используем фиксированного пользователя
	userID := uint(1)

	user, err := h.repo.User.GetUserProfile(userID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// PUT пользователя (личный кабинет)
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// В реальном приложении userID брался бы из контекста аутентификации
	userID := uint(1)

	updates := make(map[string]interface{})
	if req.Login != nil {
		updates["login"] = *req.Login
	}
	if req.Password != nil {
		updates["password"] = *req.Password
	}

	if err := h.repo.User.UpdateUserProfile(userID, updates); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// POST аутентификация
func (h *UserHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	user, err := h.repo.User.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"user_id":      user.User_ID,
		"login":        user.Login,
		"is_moderator": user.IsModerator,
	})
}

// POST деавторизация
func (h *UserHandler) Logout(ctx *gin.Context) {
	// В реальном приложении userID брался бы из контекста аутентификации
	userID := uint(1)

	if err := h.repo.User.LogoutUser(userID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
