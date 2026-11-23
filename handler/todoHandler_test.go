package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-menu/domain"
	"go-menu/handler"
	"go-menu/usecase"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newTestContext(method, path string, body []byte, params gin.Params) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	c.Request = req
	c.Params = params
	return c, w
}

func TestMenuHandler_GetAll_Success(t *testing.T) {
	port := new(mockMenuPort)
	port.On("GetAll").Return([]domain.Menu{{MenuId: 1, MenuName: "Ramen"}}, nil)
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	c, w := newTestContext(http.MethodGet, "/v1/menus", nil, nil)
	h.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp handler.MenusGetResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, []domain.Menu{{MenuId: 1, MenuName: "Ramen"}}, resp.Menus)
	port.AssertExpectations(t)
}

func TestMenuHandler_GetAll_Error(t *testing.T) {
	port := new(mockMenuPort)
	port.On("GetAll").Return(nil, errors.New("db error"))
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	c, w := newTestContext(http.MethodGet, "/v1/menus", nil, nil)
	h.GetAll(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	port.AssertExpectations(t)
}

func TestMenuHandler_CreateMenu_Success(t *testing.T) {
	port := new(mockMenuPort)
	input := domain.Menu{MenuName: "Ramen", GenreIds: []uint{1}, CategoryIds: []uint{2}}
	port.On("CreateMenu", input).Return(domain.Menu{MenuId: 1, MenuName: "Ramen", GenreIds: []uint{1}, CategoryIds: []uint{2}}, nil)
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	body, _ := json.Marshal(handler.MenuPostRequest{MenuName: "Ramen", GenreIds: []uint{1}, CategoryIds: []uint{2}})
	c, w := newTestContext(http.MethodPost, "/v1/menus", body, nil)
	h.CreateMenu(c)

	assert.Equal(t, http.StatusOK, w.Code)
	port.AssertExpectations(t)
}

func TestMenuHandler_CreateMenu_InvalidBody(t *testing.T) {
	port := new(mockMenuPort)
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	c, w := newTestContext(http.MethodPost, "/v1/menus", []byte("{invalid"), nil)
	h.CreateMenu(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	port.AssertExpectations(t)
}

func TestMenuHandler_CreateMenu_UsecaseError(t *testing.T) {
	port := new(mockMenuPort)
	input := domain.Menu{MenuName: "Ramen"}
	port.On("CreateMenu", input).Return(domain.Menu{}, errors.New("insert failed"))
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	body, _ := json.Marshal(handler.MenuPostRequest{MenuName: "Ramen"})
	c, w := newTestContext(http.MethodPost, "/v1/menus", body, nil)
	h.CreateMenu(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	port.AssertExpectations(t)
}

func TestMenuHandler_UpdateMenu_Success(t *testing.T) {
	port := new(mockMenuPort)
	input := domain.Menu{MenuId: 1, MenuName: "Ramen2"}
	port.On("UpdateMenu", input).Return(input, nil)
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	body, _ := json.Marshal(handler.MenuPutRequest{MenuName: "Ramen2"})
	c, w := newTestContext(http.MethodPut, "/v1/menus/1", body, gin.Params{{Key: "menu_id", Value: "1"}})
	h.UpdateMenu(c)

	assert.Equal(t, http.StatusOK, w.Code)
	port.AssertExpectations(t)
}

func TestMenuHandler_UpdateMenu_InvalidMenuId(t *testing.T) {
	port := new(mockMenuPort)
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	body, _ := json.Marshal(handler.MenuPutRequest{MenuName: "Ramen2"})
	c, w := newTestContext(http.MethodPut, "/v1/menus/abc", body, gin.Params{{Key: "menu_id", Value: "abc"}})
	h.UpdateMenu(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMenuHandler_UpdateMenu_InvalidBody(t *testing.T) {
	port := new(mockMenuPort)
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	c, w := newTestContext(http.MethodPut, "/v1/menus/1", []byte("{invalid"), gin.Params{{Key: "menu_id", Value: "1"}})
	h.UpdateMenu(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMenuHandler_UpdateMenu_UsecaseError(t *testing.T) {
	port := new(mockMenuPort)
	input := domain.Menu{MenuId: 1, MenuName: "Ramen2"}
	port.On("UpdateMenu", input).Return(domain.Menu{}, errors.New("update failed"))
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	body, _ := json.Marshal(handler.MenuPutRequest{MenuName: "Ramen2"})
	c, w := newTestContext(http.MethodPut, "/v1/menus/1", body, gin.Params{{Key: "menu_id", Value: "1"}})
	h.UpdateMenu(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	port.AssertExpectations(t)
}

func TestMenuHandler_UpdateGenreRelations_Success(t *testing.T) {
	port := new(mockMenuPort)
	port.On("UpdateGenreRelations", uint(1), []uint{5}).Return(domain.Menu{MenuId: 1, GenreIds: []uint{5}}, nil)
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	body, _ := json.Marshal(handler.MenuGenrePatchRequest{GenreIds: []uint{5}})
	c, w := newTestContext(http.MethodPatch, "/v1/menus/1/genres", body, gin.Params{{Key: "menu_id", Value: "1"}})
	h.UpdateGenreRelations(c)

	assert.Equal(t, http.StatusOK, w.Code)
	port.AssertExpectations(t)
}

func TestMenuHandler_UpdateGenreRelations_InvalidMenuId(t *testing.T) {
	port := new(mockMenuPort)
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	body, _ := json.Marshal(handler.MenuGenrePatchRequest{GenreIds: []uint{5}})
	c, w := newTestContext(http.MethodPatch, "/v1/menus/abc/genres", body, gin.Params{{Key: "menu_id", Value: "abc"}})
	h.UpdateGenreRelations(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMenuHandler_UpdateGenreRelations_UsecaseError(t *testing.T) {
	port := new(mockMenuPort)
	port.On("UpdateGenreRelations", uint(1), []uint{5}).Return(domain.Menu{}, errors.New("not found"))
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	body, _ := json.Marshal(handler.MenuGenrePatchRequest{GenreIds: []uint{5}})
	c, w := newTestContext(http.MethodPatch, "/v1/menus/1/genres", body, gin.Params{{Key: "menu_id", Value: "1"}})
	h.UpdateGenreRelations(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	port.AssertExpectations(t)
}

func TestMenuHandler_UpdateCategoryRelations_Success(t *testing.T) {
	port := new(mockMenuPort)
	port.On("UpdateCategoryRelations", uint(1), []uint{7}).Return(domain.Menu{MenuId: 1, CategoryIds: []uint{7}}, nil)
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	body, _ := json.Marshal(handler.MenuCategoryPatchRequest{CategoryIds: []uint{7}})
	c, w := newTestContext(http.MethodPatch, "/v1/menus/1/categories", body, gin.Params{{Key: "menu_id", Value: "1"}})
	h.UpdateCategoryRelations(c)

	assert.Equal(t, http.StatusOK, w.Code)
	port.AssertExpectations(t)
}

func TestMenuHandler_UpdateCategoryRelations_InvalidMenuId(t *testing.T) {
	port := new(mockMenuPort)
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	body, _ := json.Marshal(handler.MenuCategoryPatchRequest{CategoryIds: []uint{7}})
	c, w := newTestContext(http.MethodPatch, "/v1/menus/abc/categories", body, gin.Params{{Key: "menu_id", Value: "abc"}})
	h.UpdateCategoryRelations(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMenuHandler_UpdateCategoryRelations_UsecaseError(t *testing.T) {
	port := new(mockMenuPort)
	port.On("UpdateCategoryRelations", uint(1), []uint{7}).Return(domain.Menu{}, errors.New("not found"))
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	body, _ := json.Marshal(handler.MenuCategoryPatchRequest{CategoryIds: []uint{7}})
	c, w := newTestContext(http.MethodPatch, "/v1/menus/1/categories", body, gin.Params{{Key: "menu_id", Value: "1"}})
	h.UpdateCategoryRelations(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	port.AssertExpectations(t)
}

func TestMenuHandler_DeleteMenu_Success(t *testing.T) {
	port := new(mockMenuPort)
	port.On("DeleteMenu", uint(1)).Return(nil)
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	c, w := newTestContext(http.MethodDelete, "/v1/menus/1", nil, gin.Params{{Key: "menu_id", Value: "1"}})
	h.DeleteMenu(c)

	assert.Equal(t, http.StatusOK, w.Code)
	port.AssertExpectations(t)
}

func TestMenuHandler_DeleteMenu_InvalidMenuId(t *testing.T) {
	port := new(mockMenuPort)
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	c, w := newTestContext(http.MethodDelete, "/v1/menus/abc", nil, gin.Params{{Key: "menu_id", Value: "abc"}})
	h.DeleteMenu(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMenuHandler_DeleteMenu_UsecaseError(t *testing.T) {
	port := new(mockMenuPort)
	port.On("DeleteMenu", uint(1)).Return(errors.New("delete failed"))
	h := handler.ProvideMenuHandler(usecase.ProvideMenuUsecase(port))

	c, w := newTestContext(http.MethodDelete, "/v1/menus/1", nil, gin.Params{{Key: "menu_id", Value: "1"}})
	h.DeleteMenu(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	port.AssertExpectations(t)
}
