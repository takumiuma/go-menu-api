package user

import (
	"time"

	"gorm.io/gorm"
)

// User はAuth0統合のためのusersテーブルを表します
type User struct {
	UserID    uint   `gorm:"primaryKey;column:user_id"`
	Auth0Sub  string `gorm:"type:varchar(255);uniqueIndex;not null;column:auth0_sub"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// リレーション
	Favorites []Favorite `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// Favorite はユーザーのお気に入りメニューのためのfavoritesテーブルを表します
type Favorite struct {
	FavoriteID uint `gorm:"primaryKey;column:favorite_id"`
	UserID     uint `gorm:"not null;column:user_id;index"`
	MenuID     uint `gorm:"not null;column:menu_id;index"`
	CreatedAt  time.Time

	// リレーション
	User User `gorm:"foreignKey:UserID"`
}

func (Favorite) TableName() string {
	return "favorites"
}

// FavoriteWithMenu お気に入りとメニュー情報を含む構造体
type FavoriteWithMenu struct {
	FavoriteID  uint      `json:"favoriteId"`
	UserID      uint      `json:"userId"`
	MenuID      uint      `json:"menuId"`
	MenuName    string    `json:"menuName"`
	GenreIDs    []uint    `json:"genreIds"`
	CategoryIDs []uint    `json:"categoryIds"`
	CreatedAt   time.Time `json:"createdAt"`
}

// UserDriver はユーザー関連のデータベース操作のためのインターフェース
type UserDriver interface {
	CreateOrGetUser(auth0Sub string) (User, bool, error)
	GetUserByAuth0Sub(auth0Sub string) (User, error)
	AddFavorite(userID, menuID uint) (Favorite, error)
	RemoveFavorite(userID, menuID uint) error
	GetUserFavorites(userID uint) ([]Favorite, error)
	GetUserFavoritesWithMenu(userID uint) ([]FavoriteWithMenu, error)
}

// UserDriverImpl はUserDriverインターフェースを実装します
type UserDriverImpl struct {
	conn *gorm.DB
}

// ProvideUserDriver は新しいUserDriverImplを作成します
func ProvideUserDriver(conn *gorm.DB) UserDriver {
	return UserDriverImpl{conn: conn}
}

// CreateOrGetUser は新しいユーザーを作成するか、Auth0Subで既存のユーザーを返します
// 戻り値: (User, bool, error) - boolは新規作成の場合true
func (u UserDriverImpl) CreateOrGetUser(auth0Sub string) (User, bool, error) {
	var user User

	// 最初に既存のユーザーを検索
	err := u.conn.Where("auth0_sub = ?", auth0Sub).First(&user).Error
	if err == nil {
		// ユーザーは既に存在します
		return user, false, nil
	}

	if err != gorm.ErrRecordNotFound {
		// その他のエラーが発生しました
		return User{}, false, err
	}

	// ユーザーが存在しないため、新しいユーザーを作成
	user = User{Auth0Sub: auth0Sub}
	if err := u.conn.Create(&user).Error; err != nil {
		return User{}, false, err
	}

	return user, true, nil
}

// GetUserByAuth0Sub はAuth0 Subjectでユーザーを取得します
func (u UserDriverImpl) GetUserByAuth0Sub(auth0Sub string) (User, error) {
	var user User
	err := u.conn.Where("auth0_sub = ?", auth0Sub).First(&user).Error
	return user, err
}

// AddFavorite はメニューをユーザーのお気に入りに追加します
func (u UserDriverImpl) AddFavorite(userID, menuID uint) (Favorite, error) {
	favorite := Favorite{
		UserID: userID,
		MenuID: menuID,
	}

	err := u.conn.Create(&favorite).Error
	return favorite, err
}

// RemoveFavorite はメニューをユーザーのお気に入りから削除します
func (u UserDriverImpl) RemoveFavorite(userID, menuID uint) error {
	return u.conn.Where("user_id = ? AND menu_id = ?", userID, menuID).Delete(&Favorite{}).Error
}

// GetUserFavorites はユーザーのすべてのお気に入りを取得します
func (u UserDriverImpl) GetUserFavorites(userID uint) ([]Favorite, error) {
	var favorites []Favorite
	err := u.conn.Where("user_id = ?", userID).Find(&favorites).Error
	return favorites, err
}

// GetUserFavoritesWithMenu はユーザーのお気に入りとメニュー情報を結合して取得します
func (u UserDriverImpl) GetUserFavoritesWithMenu(userID uint) ([]FavoriteWithMenu, error) {
	var results []FavoriteWithMenu

	// favoritesテーブルとmenu_listテーブルをJOIN
	// 削除されたメニューは除外される（LEFT JOINではなくINNER JOIN）
	query := `
		SELECT 
			f.favorite_id,
			f.user_id,
			f.menu_id,
			m.menu_name,
			f.created_at
		FROM favorites f
		INNER JOIN menu_list m ON f.menu_id = m.menu_id
		WHERE f.user_id = ?
		ORDER BY f.created_at DESC
	`

	rows, err := u.conn.Raw(query, userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []struct {
		FavoriteID uint
		UserID     uint
		MenuID     uint
		MenuName   string
		CreatedAt  time.Time
	}

	for rows.Next() {
		var fav struct {
			FavoriteID uint
			UserID     uint
			MenuID     uint
			MenuName   string
			CreatedAt  time.Time
		}
		if err := u.conn.ScanRows(rows, &fav); err != nil {
			return nil, err
		}
		favorites = append(favorites, fav)
	}

	// 各お気に入りに対してジャンルとカテゴリ情報を取得
	for _, fav := range favorites {
		var genreIDs []uint
		var categoryIDs []uint

		// ジャンルIDを取得
		genreQuery := `
			SELECT mgr.genre_id 
			FROM menu_genre_relation mgr 
			WHERE mgr.menu_id = ?
		`
		if err := u.conn.Raw(genreQuery, fav.MenuID).Pluck("genre_id", &genreIDs).Error; err != nil {
			return nil, err
		}

		// カテゴリIDを取得
		categoryQuery := `
			SELECT mcr.category_id 
			FROM menu_category_relation mcr 
			WHERE mcr.menu_id = ?
		`
		if err := u.conn.Raw(categoryQuery, fav.MenuID).Pluck("category_id", &categoryIDs).Error; err != nil {
			return nil, err
		}

		result := FavoriteWithMenu{
			FavoriteID:  fav.FavoriteID,
			UserID:      fav.UserID,
			MenuID:      fav.MenuID,
			MenuName:    fav.MenuName,
			GenreIDs:    genreIDs,
			CategoryIDs: categoryIDs,
			CreatedAt:   fav.CreatedAt,
		}
		results = append(results, result)
	}

	return results, nil
}
