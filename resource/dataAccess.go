package resource

import (
	"fmt"
	"log"
	"os"

	"go-menu/resource/user"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToDatabase() *gorm.DB {
	// 環境変数から接続情報を取得
	dbUser := os.Getenv("DATASOURCE_USERNAME")
	dbPass := os.Getenv("DATASOURCE_PASSWORD")
	dbHost := os.Getenv("DATASOURCE_HOST")
	dbPort := os.Getenv("DATASOURCE_PORT")
	dbName := os.Getenv("DATASOURCE_NAME")

	// MySQL接続用のDSN（Data Source Name）を構築
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// GORMを使ってデータベースに接続
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// テーブル作成時に外部キー参照制約を無効化にすることで、マイグレーションエラーを防止
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("データベース接続に失敗しました: ", err)
	}

	// AutoMigrate実行 - UserとFavoriteモデル
	err = db.AutoMigrate(&user.User{}, &user.Favorite{})
	if err != nil {
		log.Fatal("マイグレーションに失敗しました: ", err)
	}

	// 複合ユニークインデックスの追加
	err = db.Exec("CREATE UNIQUE INDEX idx_favorites_user_menu ON favorites(user_id, menu_id)").Error
	if err != nil {
		log.Printf("複合ユニークインデックスの作成でエラーが発生しました: %v", err)
	}

	// データベース接続確認
	log.Println("データベース接続に成功しました:", db)
	log.Println("UserとFavoriteモデルのマイグレーションが完了しました")
	return db
}
