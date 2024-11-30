package port

import "go-menu/domain"

type MenuPort interface {
	GetAll() ([]domain.Menu, error)
	UpdateGenreRelations(menuId uint, genreIds []uint) (domain.Menu, error)
	UpdateCategoryRelations(menuId uint, categoryIds []uint) (domain.Menu, error)
}