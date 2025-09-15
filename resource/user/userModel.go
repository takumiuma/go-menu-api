package user

import (
	"time"

	"gorm.io/gorm"
)

// User represents the users table for Auth0 integration
type User struct {
	UserID    uint   `gorm:"primaryKey;column:user_id"`
	Auth0Sub  string `gorm:"uniqueIndex;not null;column:auth0_sub"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// リレーション
	Favorites []Favorite `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// Favorite represents the favorites table for user's favorite menus
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

// UserDriver interface for user-related database operations
type UserDriver interface {
	CreateOrGetUser(auth0Sub string) (User, error)
	GetUserByAuth0Sub(auth0Sub string) (User, error)
	AddFavorite(userID, menuID uint) (Favorite, error)
	RemoveFavorite(userID, menuID uint) error
	GetUserFavorites(userID uint) ([]Favorite, error)
}

// UserDriverImpl implements the UserDriver interface
type UserDriverImpl struct {
	conn *gorm.DB
}

// ProvideUserDriver creates a new UserDriverImpl
func ProvideUserDriver(conn *gorm.DB) UserDriver {
	return UserDriverImpl{conn: conn}
}

// CreateOrGetUser creates a new user or returns existing user by Auth0Sub
func (u UserDriverImpl) CreateOrGetUser(auth0Sub string) (User, error) {
	var user User

	// Try to find existing user first
	err := u.conn.Where("auth0_sub = ?", auth0Sub).First(&user).Error
	if err == nil {
		// User already exists
		return user, nil
	}

	if err != gorm.ErrRecordNotFound {
		// Some other error occurred
		return User{}, err
	}

	// User doesn't exist, create new one
	user = User{Auth0Sub: auth0Sub}
	if err := u.conn.Create(&user).Error; err != nil {
		return User{}, err
	}

	return user, nil
}

// GetUserByAuth0Sub retrieves a user by Auth0 Subject
func (u UserDriverImpl) GetUserByAuth0Sub(auth0Sub string) (User, error) {
	var user User
	err := u.conn.Where("auth0_sub = ?", auth0Sub).First(&user).Error
	return user, err
}

// AddFavorite adds a menu to user's favorites
func (u UserDriverImpl) AddFavorite(userID, menuID uint) (Favorite, error) {
	favorite := Favorite{
		UserID: userID,
		MenuID: menuID,
	}

	err := u.conn.Create(&favorite).Error
	return favorite, err
}

// RemoveFavorite removes a menu from user's favorites
func (u UserDriverImpl) RemoveFavorite(userID, menuID uint) error {
	return u.conn.Where("user_id = ? AND menu_id = ?", userID, menuID).Delete(&Favorite{}).Error
}

// GetUserFavorites retrieves all favorites for a user
func (u UserDriverImpl) GetUserFavorites(userID uint) ([]Favorite, error) {
	var favorites []Favorite
	err := u.conn.Where("user_id = ?", userID).Find(&favorites).Error
	return favorites, err
}
