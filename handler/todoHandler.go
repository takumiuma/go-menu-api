package handler

import (
	"go-menu/domain"
	"go-menu/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuUsecase usecase.MenuUsecase
}

func ProvideMenuHandler(u usecase.MenuUsecase) *MenuHandler {
	return &MenuHandler{u}
}

type MenuRequest struct {}

type MenusResponse struct {
	Menus []domain.Menu `json:"menus"`
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
