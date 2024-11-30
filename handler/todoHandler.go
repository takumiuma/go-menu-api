package handler

import (
	"go-menu/domain"
	"go-menu/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuUsecase usecase.MenuUsecase
}

func ProvideMenuHandler(u usecase.MenuUsecase) *MenuHandler {
	return &MenuHandler{u}
}

type MenuRequest struct{}

type MenusResponse struct {
	Menus []domain.Menu `json:"menus"`
}

type MenuGenrePatchRequest struct {
	GenreIds []uint `json:"genre_ids"`
}

type MenuCategoryPatchRequest struct {
	CategoryIds []uint `json:"category_ids"`
}

type MenuPatchResponse struct {
	Menu domain.Menu `json:"menu"`
}


func (h MenuHandler) GetAll(c *gin.Context) {
	menus, err := h.menuUsecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := MenusResponse{
		Menus: menus,
	}
	c.JSON(http.StatusOK, response)
}

func (h MenuHandler) UpdateGenreRelations(c *gin.Context) {
	var req MenuGenrePatchRequest
	// リクエストボディを取得
	if err := c.BindJSON(&req); err != nil {
		return
	}

	// パスパラメータからmenu_idを取得	
	menuId,err := strconv.Atoi(c.Param("menu_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid menu_id",
		})
		return
	}

	// ジャンルを更新
	menu, err := h.menuUsecase.UpdateGenreRelations(uint(menuId), req.GenreIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := MenuPatchResponse{
		Menu: menu,
	}

	c.JSON(http.StatusOK, response)
}

func (h MenuHandler) UpdateCategoryRelations(c *gin.Context) {
	var req MenuCategoryPatchRequest
	// リクエストボディを取得
	if err := c.BindJSON(&req); err != nil {
		return
	}

	// パスパラメータからmenu_idを取得	
	menuId,err := strconv.Atoi(c.Param("menu_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid menu_id",
		})
		return
	}

	// カテゴリを更新
	menu, err := h.menuUsecase.UpdateCategoryRelations(uint(menuId), req.CategoryIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	response := MenuPatchResponse{
		Menu: menu,
	}

	c.JSON(http.StatusOK, response)
}
