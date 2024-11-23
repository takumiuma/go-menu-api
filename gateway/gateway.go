package gateway

import (
	"go-menu/domain"
	"go-menu/resource/menu"
	"go-menu/usecase/port"
)

type MenuGateway struct {
	menuDriver menu.MenuDriver
}

func ProvideMenuPort(d menu.MenuDriver) port.MenuPort {
	return &MenuGateway{d}
}

func (t MenuGateway) GetAll() ([]domain.Menu, error) {
	result, err := t.menuDriver.GetAll()

	if err != nil {
		return nil, err
	}

	var menus []domain.Menu

	for _, t := range result {
		todo := domain.Menu{
			MenuId: t.MenuId,
			MenuName: t.MenuName,
		}
		// メニューに紐づくジャンルIDリストとカテゴリIDリストを取得
		for _, genre := range t.Genres {
			todo.GenreIds = append(todo.GenreIds, genre.GenreId)
		}
		for _, category := range t.Categories {
			todo.CategoryIds = append(todo.CategoryIds, category.CategoryId)
		}
		menus = append(menus, todo)
	}

	return menus, nil
}
