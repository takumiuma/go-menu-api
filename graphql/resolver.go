package graphql

import (
	"go-menu/domain"
	"go-menu/usecase"
)

// Resolver はGraphQLリゾルバーの基底構造体
type Resolver struct {
	MenuUsecase usecase.MenuUsecase
}

// Query はQueryリゾルバーを返す
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// Mutation はMutationリゾルバーを返す
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }

// QueryResolver インターフェース
type QueryResolver interface {
	Menus() ([]domain.Menu, error)
}

// MutationResolver インターフェース
type MutationResolver interface {
	CreateMenu(input CreateMenuInput) (*domain.Menu, error)
	UpdateMenu(input UpdateMenuInput) (*domain.Menu, error)
	DeleteMenu(menuId uint) (bool, error)
	UpdateGenreRelations(input UpdateGenreRelationsInput) (*domain.Menu, error)
	UpdateCategoryRelations(input UpdateCategoryRelationsInput) (*domain.Menu, error)
}

// Input型定義
type CreateMenuInput struct {
	MenuName    string `json:"menuName"`
	GenreIds    []uint `json:"genreIds"`
	CategoryIds []uint `json:"categoryIds"`
}

type UpdateMenuInput struct {
	MenuId      uint   `json:"menuId"`
	MenuName    string `json:"menuName"`
	GenreIds    []uint `json:"genreIds"`
	CategoryIds []uint `json:"categoryIds"`
}

type UpdateGenreRelationsInput struct {
	MenuId   uint   `json:"menuId"`
	GenreIds []uint `json:"genreIds"`
}

type UpdateCategoryRelationsInput struct {
	MenuId      uint   `json:"menuId"`
	CategoryIds []uint `json:"categoryIds"`
}

// Query実装

// Menus はメニュー一覧を取得する
func (r *queryResolver) Menus() ([]domain.Menu, error) {
	return r.MenuUsecase.GetAll()
}

// Mutation実装

// CreateMenu はメニューを作成する
func (r *mutationResolver) CreateMenu(input CreateMenuInput) (*domain.Menu, error) {
	menu := domain.Menu{
		MenuName:    input.MenuName,
		GenreIds:    input.GenreIds,
		CategoryIds: input.CategoryIds,
	}

	createdMenu, err := r.MenuUsecase.CreateMenu(menu)
	if err != nil {
		return nil, err
	}

	return &createdMenu, nil
}

// UpdateMenu はメニューを更新する
func (r *mutationResolver) UpdateMenu(input UpdateMenuInput) (*domain.Menu, error) {
	menu := domain.Menu{
		MenuId:      input.MenuId,
		MenuName:    input.MenuName,
		GenreIds:    input.GenreIds,
		CategoryIds: input.CategoryIds,
	}

	updatedMenu, err := r.MenuUsecase.UpdateMenu(menu)
	if err != nil {
		return nil, err
	}

	return &updatedMenu, nil
}

// DeleteMenu はメニューを削除する
func (r *mutationResolver) DeleteMenu(menuId uint) (bool, error) {
	err := r.MenuUsecase.DeleteMenu(menuId)
	if err != nil {
		return false, err
	}

	return true, nil
}

// UpdateGenreRelations はジャンル関連を更新する
func (r *mutationResolver) UpdateGenreRelations(input UpdateGenreRelationsInput) (*domain.Menu, error) {
	menu, err := r.MenuUsecase.UpdateGenreRelations(input.MenuId, input.GenreIds)
	if err != nil {
		return nil, err
	}

	return &menu, nil
}

// UpdateCategoryRelations はカテゴリ関連を更新する
func (r *mutationResolver) UpdateCategoryRelations(input UpdateCategoryRelationsInput) (*domain.Menu, error) {
	menu, err := r.MenuUsecase.UpdateCategoryRelations(input.MenuId, input.CategoryIds)
	if err != nil {
		return nil, err
	}

	return &menu, nil
}
