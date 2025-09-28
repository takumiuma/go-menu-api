package user

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// User はAuth0統合のためのusersテーブルを表します
type User struct {
	UserID    uint      `gorm:"primaryKey;column:user_id" json:"user_id"`
	Auth0Sub  string    `gorm:"type:varchar(255);uniqueIndex;not null;column:auth0_sub" json:"auth0_sub"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// リレーション
	Favorites []Favorite `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// Favorite はユーザーのお気に入りメニューのためのfavoritesテーブルを表します
type Favorite struct {
	FavoriteID uint      `gorm:"primaryKey;column:favorite_id" json:"favorite_id"`
	UserID     uint      `gorm:"not null;column:user_id;index" json:"user_id"`
	MenuID     uint      `gorm:"not null;column:menu_id;index" json:"menu_id"`
	CreatedAt  time.Time `json:"created_at"`

	// リレーション
	User User `gorm:"foreignKey:UserID" json:"user"`
}

func (Favorite) TableName() string {
	return "favorites"
}

// UserDriver はユーザー関連のデータベース操作のためのインターフェース
type UserDriver interface {
	CreateOrGetUser(auth0Sub string) (User, bool, error)
	GetUserByAuth0Sub(auth0Sub string) (User, error)
	AddFavorite(userID, menuID uint) (Favorite, error)
	GetUserFavorites(userID uint) ([]Favorite, error)
	GetFavoriteByID(favoriteID uint) (Favorite, error)
	RemoveFavoriteByID(favoriteID uint) error
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
	// 重複チェック：既にお気に入りに追加されているかを確認
	var existingFavorite Favorite
	err := u.conn.Where("user_id = ? AND menu_id = ?", userID, menuID).First(&existingFavorite).Error
	if err == nil {
		// 既に存在している場合は重複エラーを返す
		return Favorite{}, errors.New("Menu is already in favorites")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// その他のデータベースエラー
		return Favorite{}, err
	}

	// メニュー存在チェック：メニューテーブルにmenu_idが存在するかを確認
	var menuCount int64
	err = u.conn.Table("menu_list").Where("menu_id = ?", menuID).Count(&menuCount).Error
	if err != nil {
		return Favorite{}, err
	}
	if menuCount == 0 {
		return Favorite{}, errors.New("Menu not found")
	}

	// お気に入りを作成
	favorite := Favorite{
		UserID: userID,
		MenuID: menuID,
	}

	err = u.conn.Create(&favorite).Error
	return favorite, err
}

// GetUserFavorites はユーザーのすべてのお気に入りを取得します
func (u UserDriverImpl) GetUserFavorites(userID uint) ([]Favorite, error) {
	var favorites []Favorite
	err := u.conn.Where("user_id = ?", userID).Find(&favorites).Error
	return favorites, err
}

// GetFavoriteByID はお気に入りIDでお気に入りを取得します
func (u UserDriverImpl) GetFavoriteByID(favoriteID uint) (Favorite, error) {
	var favorite Favorite
	err := u.conn.First(&favorite, favoriteID).Error
	return favorite, err
}

// RemoveFavoriteByID はお気に入りIDでお気に入りを削除します
func (u UserDriverImpl) RemoveFavoriteByID(favoriteID uint) error {
	return u.conn.Delete(&Favorite{}, favoriteID).Error
}
