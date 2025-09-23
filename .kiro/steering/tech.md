# 技術スタック

## 言語・フレームワーク
- **Go**: 1.22.5
- **Gin**: Webフレームワーク (v1.10.0)
- **GORM**: ORMライブラリ (v1.25.7)
- **MySQL**: データベース (v1.5.7)
- **Wire**: 依存性注入 (v0.6.0)

## 主要ライブラリ
- `github.com/gin-contrib/cors`: CORS対応
- `gorm.io/driver/mysql`: MySQLドライバー
- `github.com/google/wire`: 依存性注入

## 共通コマンド

### 開発・ビルド
```bash
# アプリケーション実行
go run main.go

# ビルド
go build -o go-menu

# 依存関係の管理
go mod tidy
go mod download

# テスト実行
go test ./...
```

### サーバー設定
- デフォルトポート: `:8080`
- データベース: MySQL接続が必要

## 開発環境要件
- Go 1.22.5以上
- MySQL データベース
- 適切なデータベース接続設定