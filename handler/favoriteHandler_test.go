package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-menu/handler"
	"go-menu/resource/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func newFavoriteTestContext(method, path string, body []byte, params gin.Params, userID any, userIDSet bool) (*gin.Context, *httptest.ResponseRecorder) {
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
	if userIDSet {
		c.Set("userID", userID)
	}
	return c, w
}

func TestFavoriteHandler_AddFavorite_Unauthenticated(t *testing.T) {
	driver := new(mockUserDriver)
	h := handler.ProvideFavoriteHandler(driver)

	body, _ := json.Marshal(handler.AddFavoriteRequest{MenuID: 1})
	c, w := newFavoriteTestContext(http.MethodPost, "/v1/favorites", body, nil, nil, false)
	h.AddFavorite(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestFavoriteHandler_AddFavorite_InvalidUserIDType(t *testing.T) {
	driver := new(mockUserDriver)
	h := handler.ProvideFavoriteHandler(driver)

	body, _ := json.Marshal(handler.AddFavoriteRequest{MenuID: 1})
	c, w := newFavoriteTestContext(http.MethodPost, "/v1/favorites", body, nil, "not-a-uint", true)
	h.AddFavorite(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestFavoriteHandler_AddFavorite_InvalidBody(t *testing.T) {
	driver := new(mockUserDriver)
	h := handler.ProvideFavoriteHandler(driver)

	c, w := newFavoriteTestContext(http.MethodPost, "/v1/favorites", []byte("{invalid"), nil, uint(1), true)
	h.AddFavorite(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFavoriteHandler_AddFavorite_AlreadyExists(t *testing.T) {
	driver := new(mockUserDriver)
	driver.On("AddFavorite", uint(1), uint(2)).Return(user.Favorite{}, errors.New("Menu is already in favorites"))
	h := handler.ProvideFavoriteHandler(driver)

	body, _ := json.Marshal(handler.AddFavoriteRequest{MenuID: 2})
	c, w := newFavoriteTestContext(http.MethodPost, "/v1/favorites", body, nil, uint(1), true)
	h.AddFavorite(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	driver.AssertExpectations(t)
}

func TestFavoriteHandler_AddFavorite_MenuNotFound(t *testing.T) {
	driver := new(mockUserDriver)
	driver.On("AddFavorite", uint(1), uint(2)).Return(user.Favorite{}, errors.New("Menu not found"))
	h := handler.ProvideFavoriteHandler(driver)

	body, _ := json.Marshal(handler.AddFavoriteRequest{MenuID: 2})
	c, w := newFavoriteTestContext(http.MethodPost, "/v1/favorites", body, nil, uint(1), true)
	h.AddFavorite(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	driver.AssertExpectations(t)
}

func TestFavoriteHandler_AddFavorite_OtherError(t *testing.T) {
	driver := new(mockUserDriver)
	driver.On("AddFavorite", uint(1), uint(2)).Return(user.Favorite{}, errors.New("db error"))
	h := handler.ProvideFavoriteHandler(driver)

	body, _ := json.Marshal(handler.AddFavoriteRequest{MenuID: 2})
	c, w := newFavoriteTestContext(http.MethodPost, "/v1/favorites", body, nil, uint(1), true)
	h.AddFavorite(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	driver.AssertExpectations(t)
}

func TestFavoriteHandler_AddFavorite_Success(t *testing.T) {
	driver := new(mockUserDriver)
	driver.On("AddFavorite", uint(1), uint(2)).Return(user.Favorite{FavoriteID: 10, UserID: 1, MenuID: 2}, nil)
	h := handler.ProvideFavoriteHandler(driver)

	body, _ := json.Marshal(handler.AddFavoriteRequest{MenuID: 2})
	c, w := newFavoriteTestContext(http.MethodPost, "/v1/favorites", body, nil, uint(1), true)
	h.AddFavorite(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	driver.AssertExpectations(t)
}

func TestFavoriteHandler_GetFavorites_Unauthenticated(t *testing.T) {
	driver := new(mockUserDriver)
	h := handler.ProvideFavoriteHandler(driver)

	c, w := newFavoriteTestContext(http.MethodGet, "/v1/favorites", nil, nil, nil, false)
	h.GetFavorites(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestFavoriteHandler_GetFavorites_DriverError(t *testing.T) {
	driver := new(mockUserDriver)
	driver.On("GetUserFavorites", uint(1)).Return(nil, errors.New("db error"))
	h := handler.ProvideFavoriteHandler(driver)

	c, w := newFavoriteTestContext(http.MethodGet, "/v1/favorites", nil, nil, uint(1), true)
	h.GetFavorites(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	driver.AssertExpectations(t)
}

func TestFavoriteHandler_GetFavorites_Success(t *testing.T) {
	driver := new(mockUserDriver)
	driver.On("GetUserFavorites", uint(1)).Return([]user.Favorite{{FavoriteID: 10, UserID: 1, MenuID: 2}}, nil)
	h := handler.ProvideFavoriteHandler(driver)

	c, w := newFavoriteTestContext(http.MethodGet, "/v1/favorites", nil, nil, uint(1), true)
	h.GetFavorites(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp handler.GetFavoritesResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp.Favorites, 1)
	assert.Equal(t, uint(10), resp.Favorites[0].FavoriteID)
	driver.AssertExpectations(t)
}

func TestFavoriteHandler_RemoveFavoriteByID_Unauthenticated(t *testing.T) {
	driver := new(mockUserDriver)
	h := handler.ProvideFavoriteHandler(driver)

	c, w := newFavoriteTestContext(http.MethodDelete, "/v1/favorites/1", nil, gin.Params{{Key: "favoriteId", Value: "1"}}, nil, false)
	h.RemoveFavoriteByID(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestFavoriteHandler_RemoveFavoriteByID_InvalidID(t *testing.T) {
	driver := new(mockUserDriver)
	h := handler.ProvideFavoriteHandler(driver)

	c, w := newFavoriteTestContext(http.MethodDelete, "/v1/favorites/abc", nil, gin.Params{{Key: "favoriteId", Value: "abc"}}, uint(1), true)
	h.RemoveFavoriteByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFavoriteHandler_RemoveFavoriteByID_NotFound(t *testing.T) {
	driver := new(mockUserDriver)
	driver.On("GetFavoriteByID", uint(1)).Return(user.Favorite{}, gorm.ErrRecordNotFound)
	h := handler.ProvideFavoriteHandler(driver)

	c, w := newFavoriteTestContext(http.MethodDelete, "/v1/favorites/1", nil, gin.Params{{Key: "favoriteId", Value: "1"}}, uint(1), true)
	h.RemoveFavoriteByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	driver.AssertExpectations(t)
}

func TestFavoriteHandler_RemoveFavoriteByID_GetError(t *testing.T) {
	driver := new(mockUserDriver)
	driver.On("GetFavoriteByID", uint(1)).Return(user.Favorite{}, errors.New("db error"))
	h := handler.ProvideFavoriteHandler(driver)

	c, w := newFavoriteTestContext(http.MethodDelete, "/v1/favorites/1", nil, gin.Params{{Key: "favoriteId", Value: "1"}}, uint(1), true)
	h.RemoveFavoriteByID(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	driver.AssertExpectations(t)
}

func TestFavoriteHandler_RemoveFavoriteByID_Forbidden(t *testing.T) {
	driver := new(mockUserDriver)
	driver.On("GetFavoriteByID", uint(1)).Return(user.Favorite{FavoriteID: 1, UserID: 2}, nil)
	h := handler.ProvideFavoriteHandler(driver)

	c, w := newFavoriteTestContext(http.MethodDelete, "/v1/favorites/1", nil, gin.Params{{Key: "favoriteId", Value: "1"}}, uint(1), true)
	h.RemoveFavoriteByID(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	driver.AssertExpectations(t)
}

func TestFavoriteHandler_RemoveFavoriteByID_RemoveError(t *testing.T) {
	driver := new(mockUserDriver)
	driver.On("GetFavoriteByID", uint(1)).Return(user.Favorite{FavoriteID: 1, UserID: 1}, nil)
	driver.On("RemoveFavoriteByID", uint(1)).Return(errors.New("delete failed"))
	h := handler.ProvideFavoriteHandler(driver)

	c, w := newFavoriteTestContext(http.MethodDelete, "/v1/favorites/1", nil, gin.Params{{Key: "favoriteId", Value: "1"}}, uint(1), true)
	h.RemoveFavoriteByID(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	driver.AssertExpectations(t)
}

func TestFavoriteHandler_RemoveFavoriteByID_Success(t *testing.T) {
	driver := new(mockUserDriver)
	driver.On("GetFavoriteByID", uint(1)).Return(user.Favorite{FavoriteID: 1, UserID: 1}, nil)
	driver.On("RemoveFavoriteByID", uint(1)).Return(nil)
	h := handler.ProvideFavoriteHandler(driver)

	c, w := newFavoriteTestContext(http.MethodDelete, "/v1/favorites/1", nil, gin.Params{{Key: "favoriteId", Value: "1"}}, uint(1), true)
	h.RemoveFavoriteByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	driver.AssertExpectations(t)
}
