package handler

import (
	"go-menu/resource/menu"
	"go-menu/resource/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteHandler お気に入り機能のHTTPハンドラー
type FavoriteHandler struct {
	userDriver user.UserDriver
	menuDriver menu.MenuDriver
}

// ProvideFavoriteHandler FavoriteHandlerのコンストラクタ
func ProvideFavoriteHandler(userDriver user.UserDriver, menuDriver menu.MenuDriver) *FavoriteHandler {
	return &FavoriteHandler{userDriver: userDriver, menuDriver: menuDriver}
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

// FavoriteCheckResponse お気に入り状態確認レスポンス
type FavoriteCheckResponse struct {
	IsFavorite bool  `json:"isFavorite"`
	FavoriteID *uint `json:"favoriteId,omitempty"`
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

// RemoveFavorite お気に入りを削除
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
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

	// パスパラメータから menu_id を取得
	menuIDStr := c.Param("menu_id")
	menuID, err := strconv.ParseUint(menuIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid menu_id",
		})
		return
	}

	// お気に入りから削除
	err = h.userDriver.RemoveFavorite(userIDUint, uint(menuID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to remove favorite: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Favorite removed successfully",
	})
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

// CheckFavoriteStatus お気に入り状態を確認
func (h *FavoriteHandler) CheckFavoriteStatus(c *gin.Context) {
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

	// パスパラメータから menu_id を取得
	menuIDStr := c.Param("menuId")
	menuID, err := strconv.ParseUint(menuIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid menu_id",
		})
		return
	}

	// メニューが存在するかチェック
	exists, err = h.menuDriver.MenuExists(uint(menuID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to check menu existence: " + err.Error(),
		})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Menu not found",
		})
		return
	}

	// お気に入り状態をチェック
	isFavorite, favoriteID, err := h.userDriver.CheckFavoriteStatus(userIDUint, uint(menuID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to check favorite status: " + err.Error(),
		})
		return
	}

	response := FavoriteCheckResponse{
		IsFavorite: isFavorite,
		FavoriteID: favoriteID,
	}

	c.JSON(http.StatusOK, response)
}
