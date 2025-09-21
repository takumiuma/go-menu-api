package main

import (
	"fmt"
	"go-menu/router"
	"log"
	"os"
)

func main() {
	// 必要な環境変数の確認
	requiredEnvVars := []string{
		"DATASOURCE_USERNAME",
		"DATASOURCE_PASSWORD",
		"DATASOURCE_HOST",
		"DATASOURCE_PORT",
		"DATASOURCE_NAME",
	}

	// Auth0関連の環境変数（オプション：ない場合は警告表示）
	auth0EnvVars := []string{
		"AUTH0_DOMAIN",
		"AUTH0_AUDIENCE",
	}

	// データベース接続に必要な環境変数のチェック
	missingVars := []string{}
	for _, env := range requiredEnvVars {
		if os.Getenv(env) == "" {
			missingVars = append(missingVars, env)
		}
	}

	if len(missingVars) > 0 {
		log.Printf("警告: 以下の環境変数が設定されていません: %v", missingVars)
		log.Printf("アプリケーションの動作に影響する可能性があります")
	}

	// Auth0環境変数のチェック（警告のみ）
	missingAuth0Vars := []string{}
	for _, env := range auth0EnvVars {
		if os.Getenv(env) == "" {
			missingAuth0Vars = append(missingAuth0Vars, env)
		}
	}

	if len(missingAuth0Vars) > 0 {
		log.Printf("警告: Auth0認証に必要な環境変数が設定されていません: %v", missingAuth0Vars)
		log.Printf("認証が必要なエンドポイント(/v1/favorites/*)は動作しません")
	}

	fmt.Println("=== Go Menu API Server ===")
	fmt.Println("利用可能なエンドポイント:")
	fmt.Println("  GET    /v1/ping                           - ヘルスチェック")
	fmt.Println("  GET    /v1/menus                          - メニュー一覧取得")
	fmt.Println("  POST   /v1/menus                          - メニュー作成")
	fmt.Println("  PUT    /v1/menus/:menu_id                 - メニュー更新")
	fmt.Println("  DELETE /v1/menus/:menu_id                 - メニュー削除")
	fmt.Println("  PATCH  /v1/menus/:menu_id/genres          - ジャンル関連更新")
	fmt.Println("  PATCH  /v1/menus/:menu_id/categories       - カテゴリ関連更新")
	fmt.Println("")
	fmt.Println("認証が必要なエンドポイント (Authorization: Bearer <token>):")
	fmt.Println("  GET    /v1/favorites                      - お気に入り一覧取得")
	fmt.Println("  POST   /v1/favorites                      - お気に入り追加")
	fmt.Println("  DELETE /v1/favorites/:menu_id             - お気に入り削除")
	fmt.Println("")
	fmt.Println("サーバーを :8080 で起動しています...")

	s := router.NewServer()
	s.Run(":8080")
}
