package user

import (
	"time"

	"gorm.io/gorm"
)

// User はAuth0統合のためのusersテーブルを表します
type User struct {
	UserID    uint   `gorm:"primaryKey;column:user_id"`
	Auth0Sub  string `gorm:"uniqueIndex;not null;column:auth0_sub"`
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

// UserDriver はユーザー関連のデータベース操作のためのインターフェース
type UserDriver interface {
	CreateOrGetUser(auth0Sub string) (User, error)
	GetUserByAuth0Sub(auth0Sub string) (User, error)
	AddFavorite(userID, menuID uint) (Favorite, error)
	RemoveFavorite(userID, menuID uint) error
	GetUserFavorites(userID uint) ([]Favorite, error)
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
func (u UserDriverImpl) CreateOrGetUser(auth0Sub string) (User, error) {
	var user User

	// 最初に既存のユーザーを検索
	err := u.conn.Where("auth0_sub = ?", auth0Sub).First(&user).Error
	if err == nil {
		// ユーザーは既に存在します
		return user, nil
	}

	if err != gorm.ErrRecordNotFound {
		// その他のエラーが発生しました
		return User{}, err
	}

	// ユーザーが存在しないため、新しいユーザーを作成
	user = User{Auth0Sub: auth0Sub}
	if err := u.conn.Create(&user).Error; err != nil {
		return User{}, err
	}

	return user, nil
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
