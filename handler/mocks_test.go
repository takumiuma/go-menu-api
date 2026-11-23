package handler_test

import (
	"go-menu/domain"
	"go-menu/resource/user"

	"github.com/gin-gonic/gin"
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

type mockMenuPort struct {
	mock.Mock
}

func (m *mockMenuPort) GetAll() ([]domain.Menu, error) {
	args := m.Called()
	var result []domain.Menu
	if args.Get(0) != nil {
		result = args.Get(0).([]domain.Menu)
	}
	return result, args.Error(1)
}

func (m *mockMenuPort) CreateMenu(menu domain.Menu) (domain.Menu, error) {
	args := m.Called(menu)
	return args.Get(0).(domain.Menu), args.Error(1)
}

func (m *mockMenuPort) UpdateMenu(menu domain.Menu) (domain.Menu, error) {
	args := m.Called(menu)
	return args.Get(0).(domain.Menu), args.Error(1)
}

func (m *mockMenuPort) UpdateGenreRelations(menuId uint, genreIds []uint) (domain.Menu, error) {
	args := m.Called(menuId, genreIds)
	return args.Get(0).(domain.Menu), args.Error(1)
}

func (m *mockMenuPort) UpdateCategoryRelations(menuId uint, categoryIds []uint) (domain.Menu, error) {
	args := m.Called(menuId, categoryIds)
	return args.Get(0).(domain.Menu), args.Error(1)
}

func (m *mockMenuPort) DeleteMenu(menuId uint) error {
	args := m.Called(menuId)
	return args.Error(0)
}
