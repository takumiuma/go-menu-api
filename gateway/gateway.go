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
			EatingGenreId: t.EatingGenreId,
			EatingCategoryId: t.EatingCategoryId,
		}
		menus = append(menus, todo)
	}

	return menus, nil
}
