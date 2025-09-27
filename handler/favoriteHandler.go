package handler

import (
	"errors"
	"go-menu/resource/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FavoriteHandler お気に入り機能のHTTPハンドラー
type FavoriteHandler struct {
	userDriver user.UserDriver
}

// ProvideFavoriteHandler FavoriteHandlerのコンストラクタ
func ProvideFavoriteHandler(userDriver user.UserDriver) *FavoriteHandler {
	return &FavoriteHandler{userDriver: userDriver}
}

// AddFavoriteRequest お気に入り追加リクエスト
type AddFavoriteRequest struct {
	MenuID uint `json:"menu_id" binding:"required"`
}

// AddFavoriteResponse お気に入り追加レスポンス
type AddFavoriteResponse struct {
	Favorite user.Favorite `json:"favorite"`
}

// GetFavoritesResponse お気に入り一覧取得レスポンス
type GetFavoritesResponse struct {
	Favorites []user.Favorite `json:"favorites"`
}

// DeleteFavoriteResponse お気に入り削除レスポンス
type DeleteFavoriteResponse struct {
	Success bool `json:"success"`
}

// AddFavorite お気に入りを追加
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid user ID format",
		})
		return
	}

	// リクエストボディから menu_id を取得
	var req AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body: " + err.Error(),
		})
		return
	}

	// お気に入りに追加
	favorite, err := h.userDriver.AddFavorite(userIDUint, req.MenuID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to add favorite: " + err.Error(),
		})
		return
	}

	response := AddFavoriteResponse{
		Favorite: favorite,
	}

	c.JSON(http.StatusCreated, response)
}

// GetFavorites ユーザーのお気に入り一覧を取得
func (h *FavoriteHandler) GetFavorites(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid user ID format",
		})
		return
	}

	// ユーザーのお気に入り一覧を取得
	favorites, err := h.userDriver.GetUserFavorites(userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get favorites: " + err.Error(),
		})
		return
	}

	response := GetFavoritesResponse{
		Favorites: favorites,
	}

	c.JSON(http.StatusOK, response)
}

// RemoveFavoriteByID お気に入りをIDで削除（権限チェック付き）
func (h *FavoriteHandler) RemoveFavoriteByID(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// パスパラメータから favorite_id を取得
	favoriteIDStr := c.Param("favoriteId")
	favoriteID, err := strconv.ParseUint(favoriteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid favorite ID",
		})
		return
	}

	// お気に入りが存在するかチェック
	favorite, err := h.userDriver.GetFavoriteByID(uint(favoriteID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Favorite not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get favorite: " + err.Error(),
		})
		return
	}

	// 削除権限チェック（自分のお気に入りのみ削除可能）
	if favorite.UserID != userIDUint {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only delete your own favorites",
		})
		return
	}

	// お気に入りを削除
	err = h.userDriver.RemoveFavoriteByID(uint(favoriteID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to remove favorite: " + err.Error(),
		})
		return
	}

	response := DeleteFavoriteResponse{
		Success: true,
	}

	c.JSON(http.StatusOK, response)
}
