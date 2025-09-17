package resource

import (
	"fmt"
	"log"
	"os"

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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("データベース接続に失敗しました: ", err)
	}

	// データベース接続確認
	log.Println("データベース接続に成功しました:", db)
	return db
}
