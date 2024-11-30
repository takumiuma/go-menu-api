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
	results, err := t.menuDriver.GetAll()
	if err != nil {
		return nil, err
	}

	var menus []domain.Menu
	for _, result := range results {
		todo := domain.Menu{
			MenuId: result.MenuId,
			MenuName: result.MenuName,
			GenreIds: t.getRestGenreIds(result.Genres),
			CategoryIds: t.getRestCategoryIds(result.Categories),
		}
		menus = append(menus, todo)
	}

	return menus, nil
}

// UpdateGenreRelations はメニューに紐づくジャンルを更新する
func (t MenuGateway) UpdateGenreRelations(menuId uint, genreIds []uint) (domain.Menu, error) {
	result, err := t.menuDriver.UpdateGenreRelations(menuId, genreIds)

	if err != nil {
		return domain.Menu{}, err
	}

	menu := domain.Menu{
		MenuId: result.MenuId,
		MenuName: result.MenuName,
		GenreIds: t.getRestGenreIds(result.Genres),
		CategoryIds: t.getRestCategoryIds(result.Categories),
	}

	return menu, nil
}

// UpdateCategoryRelations はメニューに紐づくカテゴリを更新する
func (t MenuGateway) UpdateCategoryRelations(menuId uint, categoryIds []uint) (domain.Menu, error) {
	result, err := t.menuDriver.UpdateCategoryRelations(menuId, categoryIds)

	if err != nil {
		return domain.Menu{}, err
	}

	menu := domain.Menu{
		MenuId: result.MenuId,
		MenuName: result.MenuName,
		GenreIds: t.getRestGenreIds(result.Genres),
		CategoryIds: t.getRestCategoryIds(result.Categories),
	}

	return menu, nil
}

// GetRestGenreIds はメニューに紐づくジャンルIDリストを取得する
func (t MenuGateway) getRestGenreIds(genres []menu.Genre) ([]uint) {
	var genreIds []uint

	for _, genre := range genres {
		genreIds = append(genreIds, genre.GenreId)
	}

	return genreIds
}

// GetRestCategoryIds はメニューに紐づくカテゴリIDリストを取得する
func (t MenuGateway) getRestCategoryIds(categories []menu.Category) ([]uint) {
	var categoryIds []uint

	for _, category := range categories {
		categoryIds = append(categoryIds, category.CategoryId)
	}

	return categoryIds
}

