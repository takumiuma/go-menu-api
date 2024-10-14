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
	todos, err :=  u.menuPort.GetAll()

	if err != nil {
		return nil, err
	}

	return todos, nil
}
