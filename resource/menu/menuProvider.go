package menu

import (
	"gorm.io/gorm"
)

type MenuDriver interface {
	GetAll() ([]Menu, error)
	CreateMenu(menuName string, genreIds []uint, categoryIds []uint) (Menu, error)
	UpdateMenu(menuId uint, menuName string, genreIds []uint, categoryIds []uint) (Menu, error)
	UpdateGenreRelations(menuId uint, genreIds []uint) (Menu, error)
	UpdateCategoryRelations(menuId uint, categoryIds []uint) (Menu, error)
	DeleteMenu(menuId uint) error
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

// CreateMenu はメニューを作成する
func (t MenuDriverImpl) CreateMenu(menuName string, genreIds []uint, categoryIds []uint) (Menu, error) {
	menu := Menu{MenuName: menuName}

	// トランザクション開始
	tx := t.conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// メニューを作成
	if err := tx.Create(&menu).Error; err != nil {
		tx.Rollback()
		return Menu{}, err
	}

	// ジャンルを取得
	var genres []Genre
	if err := tx.Where("genre_id IN ?", genreIds).Find(&genres).Error; err != nil {
		tx.Rollback()
		return Menu{}, err
	}
	// 中間テーブルにデータを追加
	if err := tx.Model(&menu).Association("Genres").Append(genres); err != nil {
		tx.Rollback()
		return Menu{}, err
	}

	// カテゴリを取得
	var categories []Category
	if err := tx.Where("category_id IN ?", categoryIds).Find(&categories).Error; err != nil {
		tx.Rollback()
		return Menu{}, err
	}
	// 中間テーブルにデータを追加
	if err := tx.Model(&menu).Association("Categories").Append(categories); err != nil {
		tx.Rollback()
		return Menu{}, err
	}

	// コミット
	if err := tx.Commit().Error; err != nil {
		return Menu{}, err
	}

	return menu, nil
}

// UpdateMenu はメニューを更新する
func (t MenuDriverImpl) UpdateMenu(menuId uint, menuName string, genreIds []uint, categoryIds []uint) (Menu, error) {
	var menu Menu

	// メニューを取得
	if err := t.conn.Preload("Genres").Preload("Categories").First(&menu, menuId).Error; err != nil {
		return Menu{}, err
	}

	// トランザクション開始
	tx := t.conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// メニューを更新
	if err := tx.Model(&menu).Updates(Menu{MenuName: menuName}).Error; err != nil {
		tx.Rollback()
		return Menu{}, err
	}

	// ジャンルを取得
	var genres []Genre
	if err := tx.Where("genre_id IN ?", genreIds).Find(&genres).Error; err != nil {
		tx.Rollback()
		return Menu{}, err
	}
	// 中間テーブルのデータを置き換え
	if err := tx.Model(&menu).Association("Genres").Replace(genres); err != nil {
		tx.Rollback()
		return Menu{}, err
	}

	// カテゴリを取得
	var categories []Category
	if err := tx.Where("category_id IN ?", categoryIds).Find(&categories).Error; err != nil {
		tx.Rollback()
		return Menu{}, err
	}
	// 中間テーブルのデータを置き換え
	if err := tx.Model(&menu).Association("Categories").Replace(categories); err != nil {
		tx.Rollback()
		return Menu{}, err
	}

	// コミット
	if err := tx.Commit().Error; err != nil {
		return Menu{}, err
	}

	return menu, nil
}

// UpdateGenreRelations はメニューに紐づくジャンルを更新する
func (t MenuDriverImpl) UpdateGenreRelations(menuId uint, genreIds []uint) (Menu, error) {
	var menu Menu

	// メニューを取得
	if err := t.conn.Preload("Genres").Preload("Categories").First(&menu, menuId).Error; err != nil {
		return Menu{}, err
	}

	// トランザクション開始
	tx := t.conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var genres []Genre
	if err := tx.Where("genre_id IN ?", genreIds).Find(&genres).Error; err != nil {
		tx.Rollback()
		return Menu{}, err
	}
	// 中間テーブルのデータを置き換え
	if err := tx.Model(&menu).Association("Genres").Replace(genres); err != nil {
		tx.Rollback()
		return Menu{}, err
	}

	// コミット
	if err := tx.Commit().Error; err != nil {
		return Menu{}, err
	}

	return menu, nil
}

// UpdateCategoryRelations はメニューに紐づくカテゴリを更新する
func (t MenuDriverImpl) UpdateCategoryRelations(menuId uint, categoryIds []uint) (Menu, error) {
	var menu Menu

	// メニューを取得
	if err := t.conn.Preload("Genres").Preload("Categories").First(&menu, menuId).Error; err != nil {
		return Menu{}, err
	}

	// トランザクション開始
	tx := t.conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var categories []Category
	if err := tx.Where("category_id IN ?", categoryIds).Find(&categories).Error; err != nil {
		tx.Rollback()
		return Menu{}, err
	}
	// 中間テーブルのデータを置き換え
	if err := tx.Model(&menu).Association("Categories").Replace(categories); err != nil {
		tx.Rollback()
		return Menu{}, err
	}

	// コミット
	if err := tx.Commit().Error; err != nil {
		return Menu{}, err
	}

	return menu, nil
}

// DeleteMenu はメニューを削除する
func (t MenuDriverImpl) DeleteMenu(menuId uint) error {
	// トランザクション開始
	tx := t.conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// メニューを削除
	if err := tx.Delete(&Menu{}, menuId).Error; err != nil {
		tx.Rollback()
		return err
	}

	// コミット
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

type Menu struct {
	MenuId   uint   `gorm:"primaryKey" json:"id"`
	MenuName string `gorm:"size:50;column:menu_name" json:"menu_name"`
	// many2many タグで中間テーブルを指定
	Genres     []Genre    `gorm:"many2many:menu_genre_relation;joinForeignKey:menu_id;JoinReferences:genre_id"`
	Categories []Category `gorm:"many2many:menu_category_relation;joinForeignKey:menu_id;JoinReferences:category_id"`
}

func (Menu) TableName() string {
	return "menu_list"
}

type Genre struct {
	GenreId   uint   `gorm:"primaryKey" json:"genre_id"`
	GenreName string `gorm:"size:50" json:"genre_name"`
}

func (Genre) TableName() string {
	return "eating_genre_list"
}

type Category struct {
	CategoryId   uint   `gorm:"primaryKey" json:"category_id"`
	CategoryName string `gorm:"size:50" json:"category_name"`
}

func (Category) TableName() string {
	return "eating_category_list"
}
