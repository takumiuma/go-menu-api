package port

import "go-menu/domain"

type MenuPort interface {
	GetAll() ([]domain.Menu, error)
	CreateMenu(menu domain.Menu) (domain.Menu, error)
	UpdateMenu(menu domain.Menu) (domain.Menu, error)
	UpdateGenreRelations(menuId uint, genreIds []uint) (domain.Menu, error)
	UpdateCategoryRelations(menuId uint, categoryIds []uint) (domain.Menu, error)
	DeleteMenu(menuId uint) error
}
