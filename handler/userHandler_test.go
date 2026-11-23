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
)

func newUserTestContext(body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestUserHandler_CreateUser_InvalidBody(t *testing.T) {
	driver := new(mockUserDriver)
	h := handler.ProvideUserHandler(driver)

	c, w := newUserTestContext([]byte("{invalid"))
	h.CreateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_CreateUser_MissingAuth0Sub(t *testing.T) {
	driver := new(mockUserDriver)
	h := handler.ProvideUserHandler(driver)

	body, _ := json.Marshal(map[string]string{"auth0Sub": ""})
	c, w := newUserTestContext(body)
	h.CreateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_CreateUser_BlankAuth0Sub(t *testing.T) {
	driver := new(mockUserDriver)
	h := handler.ProvideUserHandler(driver)

	body, _ := json.Marshal(handler.CreateUserRequest{Auth0Sub: "   "})
	c, w := newUserTestContext(body)
	h.CreateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_CreateUser_DriverError(t *testing.T) {
	driver := new(mockUserDriver)
	driver.On("CreateOrGetUser", "auth0|123").Return(user.User{}, false, errors.New("db error"))
	h := handler.ProvideUserHandler(driver)

	body, _ := json.Marshal(handler.CreateUserRequest{Auth0Sub: "auth0|123"})
	c, w := newUserTestContext(body)
	h.CreateUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	driver.AssertExpectations(t)
}

func TestUserHandler_CreateUser_NewUser(t *testing.T) {
	driver := new(mockUserDriver)
	driver.On("CreateOrGetUser", "auth0|123").Return(user.User{UserID: 1, Auth0Sub: "auth0|123"}, true, nil)
	h := handler.ProvideUserHandler(driver)

	body, _ := json.Marshal(handler.CreateUserRequest{Auth0Sub: "auth0|123"})
	c, w := newUserTestContext(body)
	h.CreateUser(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	driver.AssertExpectations(t)
}

func TestUserHandler_CreateUser_ExistingUser(t *testing.T) {
	driver := new(mockUserDriver)
	driver.On("CreateOrGetUser", "auth0|123").Return(user.User{UserID: 1, Auth0Sub: "auth0|123"}, false, nil)
	h := handler.ProvideUserHandler(driver)

	body, _ := json.Marshal(handler.CreateUserRequest{Auth0Sub: "auth0|123"})
	c, w := newUserTestContext(body)
	h.CreateUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	driver.AssertExpectations(t)
}
