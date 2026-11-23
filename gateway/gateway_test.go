package gateway_test

import (
	"errors"
	"testing"

	"go-menu/domain"
	"go-menu/gateway"
	"go-menu/resource/menu"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockMenuDriver struct {
	mock.Mock
}

func (m *mockMenuDriver) GetAll() ([]menu.Menu, error) {
	args := m.Called()
	var result []menu.Menu
	if args.Get(0) != nil {
		result = args.Get(0).([]menu.Menu)
	}
	return result, args.Error(1)
}

func (m *mockMenuDriver) CreateMenu(menuName string, genreIds []uint, categoryIds []uint) (menu.Menu, error) {
	args := m.Called(menuName, genreIds, categoryIds)
	return args.Get(0).(menu.Menu), args.Error(1)
}

func (m *mockMenuDriver) UpdateMenu(menuId uint, menuName string, genreIds []uint, categoryIds []uint) (menu.Menu, error) {
	args := m.Called(menuId, menuName, genreIds, categoryIds)
	return args.Get(0).(menu.Menu), args.Error(1)
}

func (m *mockMenuDriver) UpdateGenreRelations(menuId uint, genreIds []uint) (menu.Menu, error) {
	args := m.Called(menuId, genreIds)
	return args.Get(0).(menu.Menu), args.Error(1)
}

func (m *mockMenuDriver) UpdateCategoryRelations(menuId uint, categoryIds []uint) (menu.Menu, error) {
	args := m.Called(menuId, categoryIds)
	return args.Get(0).(menu.Menu), args.Error(1)
}

func (m *mockMenuDriver) DeleteMenu(menuId uint) error {
	args := m.Called(menuId)
	return args.Error(0)
}

func TestMenuGateway_GetAll_Success(t *testing.T) {
	driver := new(mockMenuDriver)
	driver.On("GetAll").Return([]menu.Menu{
		{
			MenuId:     1,
			MenuName:   "Ramen",
			Genres:     []menu.Genre{{GenreId: 10}, {GenreId: 20}},
			Categories: []menu.Category{{CategoryId: 100}},
		},
	}, nil)

	port := gateway.ProvideMenuPort(driver)
	result, err := port.GetAll()

	assert.NoError(t, err)
	assert.Equal(t, []domain.Menu{
		{
			MenuId:      1,
			MenuName:    "Ramen",
			GenreIds:    []uint{10, 20},
			CategoryIds: []uint{100},
		},
	}, result)
	driver.AssertExpectations(t)
}

func TestMenuGateway_GetAll_Error(t *testing.T) {
	driver := new(mockMenuDriver)
	driver.On("GetAll").Return(nil, errors.New("db error"))

	port := gateway.ProvideMenuPort(driver)
	result, err := port.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	driver.AssertExpectations(t)
}

func TestMenuGateway_CreateMenu_Success(t *testing.T) {
	driver := new(mockMenuDriver)
	driver.On("CreateMenu", "Ramen", []uint{10}, []uint{100}).Return(menu.Menu{
		MenuId:     1,
		MenuName:   "Ramen",
		Genres:     []menu.Genre{{GenreId: 10}},
		Categories: []menu.Category{{CategoryId: 100}},
	}, nil)

	port := gateway.ProvideMenuPort(driver)
	result, err := port.CreateMenu(domain.Menu{MenuName: "Ramen", GenreIds: []uint{10}, CategoryIds: []uint{100}})

	assert.NoError(t, err)
	assert.Equal(t, domain.Menu{MenuId: 1, MenuName: "Ramen", GenreIds: []uint{10}, CategoryIds: []uint{100}}, result)
	driver.AssertExpectations(t)
}

func TestMenuGateway_CreateMenu_Error(t *testing.T) {
	driver := new(mockMenuDriver)
	driver.On("CreateMenu", "Ramen", []uint(nil), []uint(nil)).Return(menu.Menu{}, errors.New("insert failed"))

	port := gateway.ProvideMenuPort(driver)
	result, err := port.CreateMenu(domain.Menu{MenuName: "Ramen"})

	assert.Error(t, err)
	assert.Equal(t, domain.Menu{}, result)
	driver.AssertExpectations(t)
}

func TestMenuGateway_UpdateMenu_Success(t *testing.T) {
	driver := new(mockMenuDriver)
	driver.On("UpdateMenu", uint(1), "Ramen2", []uint{20}, []uint{200}).Return(menu.Menu{
		MenuId:     1,
		MenuName:   "Ramen2",
		Genres:     []menu.Genre{{GenreId: 20}},
		Categories: []menu.Category{{CategoryId: 200}},
	}, nil)

	port := gateway.ProvideMenuPort(driver)
	result, err := port.UpdateMenu(domain.Menu{MenuId: 1, MenuName: "Ramen2", GenreIds: []uint{20}, CategoryIds: []uint{200}})

	assert.NoError(t, err)
	assert.Equal(t, domain.Menu{MenuId: 1, MenuName: "Ramen2", GenreIds: []uint{20}, CategoryIds: []uint{200}}, result)
	driver.AssertExpectations(t)
}

func TestMenuGateway_UpdateMenu_Error(t *testing.T) {
	driver := new(mockMenuDriver)
	driver.On("UpdateMenu", uint(1), "Ramen2", []uint(nil), []uint(nil)).Return(menu.Menu{}, errors.New("update failed"))

	port := gateway.ProvideMenuPort(driver)
	result, err := port.UpdateMenu(domain.Menu{MenuId: 1, MenuName: "Ramen2"})

	assert.Error(t, err)
	assert.Equal(t, domain.Menu{}, result)
	driver.AssertExpectations(t)
}

func TestMenuGateway_UpdateGenreRelations_Success(t *testing.T) {
	driver := new(mockMenuDriver)
	driver.On("UpdateGenreRelations", uint(1), []uint{30}).Return(menu.Menu{
		MenuId:   1,
		MenuName: "Ramen",
		Genres:   []menu.Genre{{GenreId: 30}},
	}, nil)

	port := gateway.ProvideMenuPort(driver)
	result, err := port.UpdateGenreRelations(1, []uint{30})

	assert.NoError(t, err)
	assert.Equal(t, domain.Menu{MenuId: 1, MenuName: "Ramen", GenreIds: []uint{30}}, result)
	driver.AssertExpectations(t)
}

func TestMenuGateway_UpdateGenreRelations_Error(t *testing.T) {
	driver := new(mockMenuDriver)
	driver.On("UpdateGenreRelations", uint(1), []uint{30}).Return(menu.Menu{}, errors.New("not found"))

	port := gateway.ProvideMenuPort(driver)
	result, err := port.UpdateGenreRelations(1, []uint{30})

	assert.Error(t, err)
	assert.Equal(t, domain.Menu{}, result)
	driver.AssertExpectations(t)
}

func TestMenuGateway_UpdateCategoryRelations_Success(t *testing.T) {
	driver := new(mockMenuDriver)
	driver.On("UpdateCategoryRelations", uint(1), []uint{40}).Return(menu.Menu{
		MenuId:     1,
		MenuName:   "Ramen",
		Categories: []menu.Category{{CategoryId: 40}},
	}, nil)

	port := gateway.ProvideMenuPort(driver)
	result, err := port.UpdateCategoryRelations(1, []uint{40})

	assert.NoError(t, err)
	assert.Equal(t, domain.Menu{MenuId: 1, MenuName: "Ramen", CategoryIds: []uint{40}}, result)
	driver.AssertExpectations(t)
}

func TestMenuGateway_UpdateCategoryRelations_Error(t *testing.T) {
	driver := new(mockMenuDriver)
	driver.On("UpdateCategoryRelations", uint(1), []uint{40}).Return(menu.Menu{}, errors.New("not found"))

	port := gateway.ProvideMenuPort(driver)
	result, err := port.UpdateCategoryRelations(1, []uint{40})

	assert.Error(t, err)
	assert.Equal(t, domain.Menu{}, result)
	driver.AssertExpectations(t)
}

func TestMenuGateway_DeleteMenu_Success(t *testing.T) {
	driver := new(mockMenuDriver)
	driver.On("DeleteMenu", uint(1)).Return(nil)

	port := gateway.ProvideMenuPort(driver)
	err := port.DeleteMenu(1)

	assert.NoError(t, err)
	driver.AssertExpectations(t)
}

func TestMenuGateway_DeleteMenu_Error(t *testing.T) {
	driver := new(mockMenuDriver)
	driver.On("DeleteMenu", uint(1)).Return(errors.New("delete failed"))

	port := gateway.ProvideMenuPort(driver)
	err := port.DeleteMenu(1)

	assert.Error(t, err)
	driver.AssertExpectations(t)
}
