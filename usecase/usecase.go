package usecase

import (
	"go-menu/domain"
	"go-menu/usecase/port"
)

type MenuUsecase struct {
	menuPort port.MenuPort
}

func ProvideMenuUsecase(menuPort port.MenuPort) MenuUsecase {
	return MenuUsecase{menuPort}
}

func (u MenuUsecase) GetAll() ([]domain.Menu, error) {
	menus, err :=  u.menuPort.GetAll()

	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (u MenuUsecase) CreateMenu(menu domain.Menu) (domain.Menu, error) {
    menus, err := u.menuPort.CreateMenu(menu)
	if err != nil {
		return domain.Menu{}, err
	}

	return menus, nil
}

func (u MenuUsecase) UpdateMenu(menu domain.Menu) (domain.Menu, error) {
	menu, err := u.menuPort.UpdateMenu(menu)

	if err != nil {
		return domain.Menu{}, err
	}

	return menu, nil
}

func (u MenuUsecase) UpdateGenreRelations(menuId uint, genreIds []uint) (domain.Menu, error) {
	menu, err := u.menuPort.UpdateGenreRelations(menuId, genreIds)

	if err != nil {
		return domain.Menu{}, err
	}

	return menu, nil
}

func (u MenuUsecase) UpdateCategoryRelations(menuId uint, categoryIds []uint) (domain.Menu, error) {
	menu, err := u.menuPort.UpdateCategoryRelations(menuId, categoryIds)

	if err != nil {
		return domain.Menu{}, err
	}

	return menu, nil
}

func (u MenuUsecase) DeleteMenu(menuId uint) error {
	err := u.menuPort.DeleteMenu(menuId)

	if err != nil {
		return err
	}

	return nil
}
