package usecase_test

import (
	"errors"
	"testing"

	"go-menu/domain"
	"go-menu/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestMenuUsecase_GetAll_Success(t *testing.T) {
	port := new(mockMenuPort)
	expected := []domain.Menu{{MenuId: 1, MenuName: "Ramen"}}
	port.On("GetAll").Return(expected, nil)

	u := usecase.ProvideMenuUsecase(port)
	result, err := u.GetAll()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	port.AssertExpectations(t)
}

func TestMenuUsecase_GetAll_Error(t *testing.T) {
	port := new(mockMenuPort)
	port.On("GetAll").Return(nil, errors.New("db error"))

	u := usecase.ProvideMenuUsecase(port)
	result, err := u.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	port.AssertExpectations(t)
}

func TestMenuUsecase_CreateMenu_Success(t *testing.T) {
	port := new(mockMenuPort)
	input := domain.Menu{MenuName: "Ramen"}
	expected := domain.Menu{MenuId: 1, MenuName: "Ramen"}
	port.On("CreateMenu", input).Return(expected, nil)

	u := usecase.ProvideMenuUsecase(port)
	result, err := u.CreateMenu(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	port.AssertExpectations(t)
}

func TestMenuUsecase_CreateMenu_Error(t *testing.T) {
	port := new(mockMenuPort)
	input := domain.Menu{MenuName: "Ramen"}
	port.On("CreateMenu", input).Return(domain.Menu{}, errors.New("insert failed"))

	u := usecase.ProvideMenuUsecase(port)
	result, err := u.CreateMenu(input)

	assert.Error(t, err)
	assert.Equal(t, domain.Menu{}, result)
	port.AssertExpectations(t)
}

func TestMenuUsecase_UpdateMenu_Success(t *testing.T) {
	port := new(mockMenuPort)
	input := domain.Menu{MenuId: 1, MenuName: "Ramen2"}
	port.On("UpdateMenu", input).Return(input, nil)

	u := usecase.ProvideMenuUsecase(port)
	result, err := u.UpdateMenu(input)

	assert.NoError(t, err)
	assert.Equal(t, input, result)
	port.AssertExpectations(t)
}

func TestMenuUsecase_UpdateMenu_Error(t *testing.T) {
	port := new(mockMenuPort)
	input := domain.Menu{MenuId: 1, MenuName: "Ramen2"}
	port.On("UpdateMenu", input).Return(domain.Menu{}, errors.New("update failed"))

	u := usecase.ProvideMenuUsecase(port)
	result, err := u.UpdateMenu(input)

	assert.Error(t, err)
	assert.Equal(t, domain.Menu{}, result)
	port.AssertExpectations(t)
}

func TestMenuUsecase_UpdateGenreRelations_Success(t *testing.T) {
	port := new(mockMenuPort)
	expected := domain.Menu{MenuId: 1, GenreIds: []uint{10}}
	port.On("UpdateGenreRelations", uint(1), []uint{10}).Return(expected, nil)

	u := usecase.ProvideMenuUsecase(port)
	result, err := u.UpdateGenreRelations(1, []uint{10})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	port.AssertExpectations(t)
}

func TestMenuUsecase_UpdateGenreRelations_Error(t *testing.T) {
	port := new(mockMenuPort)
	port.On("UpdateGenreRelations", uint(1), []uint{10}).Return(domain.Menu{}, errors.New("not found"))

	u := usecase.ProvideMenuUsecase(port)
	result, err := u.UpdateGenreRelations(1, []uint{10})

	assert.Error(t, err)
	assert.Equal(t, domain.Menu{}, result)
	port.AssertExpectations(t)
}

func TestMenuUsecase_UpdateCategoryRelations_Success(t *testing.T) {
	port := new(mockMenuPort)
	expected := domain.Menu{MenuId: 1, CategoryIds: []uint{20}}
	port.On("UpdateCategoryRelations", uint(1), []uint{20}).Return(expected, nil)

	u := usecase.ProvideMenuUsecase(port)
	result, err := u.UpdateCategoryRelations(1, []uint{20})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	port.AssertExpectations(t)
}

func TestMenuUsecase_UpdateCategoryRelations_Error(t *testing.T) {
	port := new(mockMenuPort)
	port.On("UpdateCategoryRelations", uint(1), []uint{20}).Return(domain.Menu{}, errors.New("not found"))

	u := usecase.ProvideMenuUsecase(port)
	result, err := u.UpdateCategoryRelations(1, []uint{20})

	assert.Error(t, err)
	assert.Equal(t, domain.Menu{}, result)
	port.AssertExpectations(t)
}

func TestMenuUsecase_DeleteMenu_Success(t *testing.T) {
	port := new(mockMenuPort)
	port.On("DeleteMenu", uint(1)).Return(nil)

	u := usecase.ProvideMenuUsecase(port)
	err := u.DeleteMenu(1)

	assert.NoError(t, err)
	port.AssertExpectations(t)
}

func TestMenuUsecase_DeleteMenu_Error(t *testing.T) {
	port := new(mockMenuPort)
	port.On("DeleteMenu", uint(1)).Return(errors.New("delete failed"))

	u := usecase.ProvideMenuUsecase(port)
	err := u.DeleteMenu(1)

	assert.Error(t, err)
	port.AssertExpectations(t)
}
