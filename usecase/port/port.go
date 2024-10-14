package port

import "go-menu/domain"

type MenuPort interface {
	GetAll() ([]domain.Menu, error)
}