package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-menu/handler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSystemHandler_Ping(t *testing.T) {
	h := handler.NewSystemHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/ping", nil)

	h.Ping(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message":"pong","status":"healthy"}`, w.Body.String())
}
