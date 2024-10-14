package menu

import (
	"gorm.io/gorm"
)

type MenuDriver interface {
	GetAll() ([]Menu, error)
}

type MenuDriverImpl struct {
	conn *gorm.DB
}

func ProvideMenuDriver(conn *gorm.DB) MenuDriver {
	return MenuDriverImpl{conn: conn}
}

func (t MenuDriverImpl) GetAll() ([]Menu, error) {
	menus := []Menu{}
	t.conn.Find(&menus)

	return menus, nil
}

type Menu struct{
	MenuId			uint   `gorm:"primaryKey" json:"id"`
	MenuName	string `gorm:"size:50" json:"menu_name"`
	EatingGenreId uint `json:"eating_genre_id"`
	EatingCategoryId uint `json:"eating_category_id"`
}

func (Menu) TableName() string {
	return "menu_list"
}

type Genre struct{
	EatingGenreId			uint   `gorm:"primaryKey" json:"genre_id"`
	GenreName	string `gorm:"size:50" json:"genre_name"`
}

func (Genre) TableName() string {
	return "eating_genre_list"
}

type Categories struct{
	EatingCategoryId			uint   `gorm:"primaryKey" json:"category_id"`
	CategoryName	string `gorm:"size:50" json:"category_name"`
}

func (Categories) TableName() string {
	return "eating_category_list"
}
