package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-menu/resource/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type mockUserDriver struct {
	mock.Mock
}

func (m *mockUserDriver) CreateOrGetUser(auth0Sub string) (user.User, bool, error) {
	args := m.Called(auth0Sub)
	return args.Get(0).(user.User), args.Bool(1), args.Error(2)
}

func (m *mockUserDriver) GetUserByAuth0Sub(auth0Sub string) (user.User, error) {
	args := m.Called(auth0Sub)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *mockUserDriver) AddFavorite(userID, menuID uint) (user.Favorite, error) {
	args := m.Called(userID, menuID)
	return args.Get(0).(user.Favorite), args.Error(1)
}

func (m *mockUserDriver) GetUserFavorites(userID uint) ([]user.Favorite, error) {
	args := m.Called(userID)
	var result []user.Favorite
	if args.Get(0) != nil {
		result = args.Get(0).([]user.Favorite)
	}
	return result, args.Error(1)
}

func (m *mockUserDriver) GetFavoriteByID(favoriteID uint) (user.Favorite, error) {
	args := m.Called(favoriteID)
	return args.Get(0).(user.Favorite), args.Error(1)
}

func (m *mockUserDriver) RemoveFavoriteByID(favoriteID uint) error {
	args := m.Called(favoriteID)
	return args.Error(0)
}

func setupAuthRouter(driver user.UserDriver) *gin.Engine {
	r := gin.New()
	r.Use(AuthMiddleware(driver, Auth0Config{Domain: "example.auth0.com", Audience: "test-audience"}))
	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return r
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	driver := new(mockUserDriver)
	r := setupAuthRouter(driver)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	driver.AssertNotCalled(t, "GetUserByAuth0Sub")
}

func TestAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	driver := new(mockUserDriver)
	r := setupAuthRouter(driver)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Basic abcdef")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	driver.AssertNotCalled(t, "GetUserByAuth0Sub")
}

func TestAuthMiddleware_MalformedToken(t *testing.T) {
	driver := new(mockUserDriver)
	r := setupAuthRouter(driver)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-jwt")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	driver.AssertNotCalled(t, "GetUserByAuth0Sub")
}

func TestNewAuth0Config(t *testing.T) {
	t.Setenv("AUTH0_DOMAIN", "example.auth0.com")
	t.Setenv("AUTH0_AUDIENCE", "test-audience")

	config := NewAuth0Config()

	assert.Equal(t, Auth0Config{Domain: "example.auth0.com", Audience: "test-audience"}, config)
}

func TestNewAuth0Config_Empty(t *testing.T) {
	t.Setenv("AUTH0_DOMAIN", "")
	t.Setenv("AUTH0_AUDIENCE", "")

	config := NewAuth0Config()

	assert.Equal(t, Auth0Config{}, config)
}

func TestConvertJWKToRSAPublicKey_Success(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	eBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(eBytes, uint32(privateKey.E))
	// 先頭のゼロバイトを除去（JWKのeは最小バイト数で表現される）
	for len(eBytes) > 1 && eBytes[0] == 0 {
		eBytes = eBytes[1:]
	}

	jwk := JWK{
		Kty: "RSA",
		Use: "sig",
		Kid: "test-kid",
		N:   base64.RawURLEncoding.EncodeToString(privateKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(eBytes),
	}

	publicKey, err := convertJWKToRSAPublicKey(jwk)

	assert.NoError(t, err)
	assert.Equal(t, privateKey.E, publicKey.E)
	assert.Equal(t, 0, privateKey.N.Cmp(publicKey.N))
}

func TestConvertJWKToRSAPublicKey_InvalidN(t *testing.T) {
	jwk := JWK{
		N: "not-valid-base64!!!",
		E: base64.RawURLEncoding.EncodeToString([]byte{0x01, 0x00, 0x01}),
	}

	_, err := convertJWKToRSAPublicKey(jwk)

	assert.Error(t, err)
}

func TestConvertJWKToRSAPublicKey_InvalidE(t *testing.T) {
	jwk := JWK{
		N: base64.RawURLEncoding.EncodeToString([]byte{0x01, 0x02, 0x03}),
		E: "not-valid-base64!!!",
	}

	_, err := convertJWKToRSAPublicKey(jwk)

	assert.Error(t, err)
}
