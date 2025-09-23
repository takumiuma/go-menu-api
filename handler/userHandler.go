package handler

import (
	"go-menu/resource/user"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// UserHandler ユーザー管理のHTTPハンドラー
type UserHandler struct {
	userDriver user.UserDriver
}

// ProvideUserHandler UserHandlerのコンストラクタ
func ProvideUserHandler(userDriver user.UserDriver) *UserHandler {
	return &UserHandler{userDriver: userDriver}
}

// CreateUserRequest ユーザー作成リクエスト
type CreateUserRequest struct {
	Auth0Sub string `json:"auth0Sub" binding:"required"`
}

// UserResponse ユーザーレスポンス
type UserResponse struct {
	User user.User `json:"user"`
}

// CreateUser ユーザーを作成または取得する
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "リクエストボディが無効です: " + err.Error(),
		})
		return
	}

	// Auth0 sub のフォーマットを簡単に検証
	if strings.TrimSpace(req.Auth0Sub) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "auth0Sub は必須です",
		})
		return
	}

	// ユーザーを作成または取得
	userRecord, isNewUser, err := h.userDriver.CreateOrGetUser(req.Auth0Sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ユーザーの作成または取得に失敗しました: " + err.Error(),
		})
		return
	}

	response := UserResponse{
		User: userRecord,
	}

	// 新規作成時は201 Created、既存ユーザー時は200 OK
	if isNewUser {
		c.JSON(http.StatusCreated, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}
