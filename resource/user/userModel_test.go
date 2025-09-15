package user

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// AutoMigrate the schema
	db.AutoMigrate(&User{}, &Favorite{})

	// Create the composite unique index
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_favorites_user_menu ON favorites(user_id, menu_id)")

	return db
}

func TestUserModel(t *testing.T) {
	db := setupTestDB()

	// Test User model creation
	user := User{
		Auth0Sub: "auth0|test123456789",
	}

	if err := db.Create(&user).Error; err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	// Test user was created with auto-generated ID
	if user.UserID == 0 {
		t.Error("User ID should be auto-generated")
	}

	// Test Auth0Sub unique constraint
	duplicateUser := User{
		Auth0Sub: "auth0|test123456789",
	}

	if err := db.Create(&duplicateUser).Error; err == nil {
		t.Error("Should not allow duplicate Auth0Sub")
	}
}

func TestFavoriteModel(t *testing.T) {
	db := setupTestDB()

	// Create a test user first
	user := User{
		Auth0Sub: "auth0|test123456789",
	}
	db.Create(&user)

	// Test Favorite model creation
	favorite := Favorite{
		UserID: user.UserID,
		MenuID: 1,
	}

	if err := db.Create(&favorite).Error; err != nil {
		t.Errorf("Failed to create favorite: %v", err)
	}

	// Test composite unique constraint
	duplicateFavorite := Favorite{
		UserID: user.UserID,
		MenuID: 1,
	}

	if err := db.Create(&duplicateFavorite).Error; err == nil {
		t.Error("Should not allow duplicate user_id and menu_id combination")
	}
}

func TestUserDriverCreateOrGetUser(t *testing.T) {
	db := setupTestDB()
	driver := ProvideUserDriver(db)

	auth0Sub := "auth0|test123456789"

	// Test creating new user
	user1, err := driver.CreateOrGetUser(auth0Sub)
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	if user1.Auth0Sub != auth0Sub {
		t.Errorf("Expected Auth0Sub %s, got %s", auth0Sub, user1.Auth0Sub)
	}

	// Test getting existing user
	user2, err := driver.CreateOrGetUser(auth0Sub)
	if err != nil {
		t.Errorf("Failed to get existing user: %v", err)
	}

	if user1.UserID != user2.UserID {
		t.Error("Should return the same user for existing Auth0Sub")
	}
}

func TestUserDriverFavorites(t *testing.T) {
	db := setupTestDB()
	driver := ProvideUserDriver(db)

	// Create a test user
	user, err := driver.CreateOrGetUser("auth0|test123456789")
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	menuID := uint(1)

	// Test adding favorite
	favorite, err := driver.AddFavorite(user.UserID, menuID)
	if err != nil {
		t.Errorf("Failed to add favorite: %v", err)
	}

	if favorite.UserID != user.UserID || favorite.MenuID != menuID {
		t.Error("Favorite not created correctly")
	}

	// Test getting user favorites
	favorites, err := driver.GetUserFavorites(user.UserID)
	if err != nil {
		t.Errorf("Failed to get user favorites: %v", err)
	}

	if len(favorites) != 1 {
		t.Errorf("Expected 1 favorite, got %d", len(favorites))
	}

	// Test removing favorite
	err = driver.RemoveFavorite(user.UserID, menuID)
	if err != nil {
		t.Errorf("Failed to remove favorite: %v", err)
	}

	// Verify favorite was removed
	favorites, err = driver.GetUserFavorites(user.UserID)
	if err != nil {
		t.Errorf("Failed to get user favorites after removal: %v", err)
	}

	if len(favorites) != 0 {
		t.Errorf("Expected 0 favorites after removal, got %d", len(favorites))
	}
}