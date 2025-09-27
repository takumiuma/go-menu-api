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
	favorite := Favorite{
		UserID: userID,
		MenuID: menuID,
	}

	err := u.conn.Create(&favorite).Error
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
