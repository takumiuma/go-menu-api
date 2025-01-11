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

type MenusGetResponse struct {
	Menus []domain.Menu `json:"menus"`
}

type MenuPostRequest struct {
    MenuName   string `json:"menu_name"`
    GenreIds   []uint `json:"genre_ids"`
    CategoryIds []uint `json:"category_ids"`
}

type MenuPostResponse struct {
    Menu domain.Menu `json:"menu"`
}

type MenuPutRequest struct {
	MenuName   string `json:"menu_name"`
    GenreIds   []uint `json:"genre_ids"`
    CategoryIds []uint `json:"category_ids"`
}

type MenuPutResponse struct {
	Menu domain.Menu `json:"menu"`
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

	response := MenusGetResponse{
		Menus: menus,
	}
	c.JSON(http.StatusOK, response)
}

func (h MenuHandler) CreateMenu(c *gin.Context) {
    var req MenuPostRequest
    // リクエストボディを取得
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request",
        })
        return
    }

    // メニューを作成
    menu := domain.Menu{
        MenuName:   req.MenuName,
        GenreIds:   req.GenreIds,
        CategoryIds: req.CategoryIds,
    }

    createdMenu, err := h.menuUsecase.CreateMenu(menu)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": err.Error(),
        })
        return
    }

    response := MenuPostResponse{
        Menu: createdMenu,
    }

    c.JSON(http.StatusOK, response)
}

func (h MenuHandler) UpdateMenu(c *gin.Context) {
	var req MenuPutRequest
	// リクエストボディを取得
	if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "invalid request",
        })
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

	// メニューを更新
	menu := domain.Menu{
		MenuId: uint(menuId),
		MenuName: req.MenuName,
		GenreIds: req.GenreIds,
		CategoryIds: req.CategoryIds,
	}

	updatedMenu, err := h.menuUsecase.UpdateMenu(menu)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": err.Error(),
        })
        return
    }

    response := MenuPutResponse{
        Menu: updatedMenu,
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
