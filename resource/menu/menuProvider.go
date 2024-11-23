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
	// Preloadで関連データを読み込む
	// 動作が遅くなる場合はPluckかJoinsを使って最適化する
    if err := t.conn.Preload("Genres").Preload("Categories").Find(&menus).Error; err != nil {
        return nil, err
    }

	return menus, nil
}

type Menu struct{
	MenuId			uint   `gorm:"primaryKey" json:"id"`
	MenuName	string `gorm:"size:50;column:menu_name" json:"menu_name"`
	// many2many タグで中間テーブルを指定
    Genres      []Genre         `gorm:"many2many:menu_genre_relation;joinForeignKey:menu_id;JoinReferences:genre_id"`
    Categories  []Category      `gorm:"many2many:menu_category_relation;joinForeignKey:menu_id;JoinReferences:category_id"`
}

func (Menu) TableName() string {
	return "menu_list"
}

type Genre struct{
	GenreId			uint   `gorm:"primaryKey" json:"genre_id"`
	GenreName	string `gorm:"size:50" json:"genre_name"`
}

func (Genre) TableName() string {
	return "eating_genre_list"
}

type Category struct{
	CategoryId			uint   `gorm:"primaryKey" json:"category_id"`
	CategoryName	string `gorm:"size:50" json:"category_name"`
}

func (Category) TableName() string {
	return "eating_category_list"
}
