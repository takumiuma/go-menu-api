# Auth0認証ミドルウェア実装ガイド

## 概要

このドキュメントでは、Go Menu APIに実装されたAuth0 JWTトークン検証ミドルウェアの使い方と設定方法について説明します。

## 機能

### 実装された機能

- ✅ Authorization Bearerヘッダーからトークンを正しく抽出
- ✅ Auth0の公開鍵でJWTトークンを検証
- ✅ JWTトークンからAuth0 subを抽出
- ✅ Auth0 subからデータベースのユーザーIDを取得
- ✅ 認証エラー時に適切なHTTPステータスコードを返す
- ✅ ミドルウェアがGinフレームワークで動作

### 新しく追加されたAPIエンドポイント

#### 認証が必要なお気に入り関連API

| Method | Endpoint | 説明 | 認証 |
|--------|----------|------|------|
| GET | `/v1/favorites` | ユーザーのお気に入り一覧取得 | ✅ 必要 |
| POST | `/v1/favorites` | お気に入りメニュー追加 | ✅ 必要 |
| DELETE | `/v1/favorites/:menu_id` | お気に入りメニュー削除 | ✅ 必要 |

## 環境変数設定

### 必須環境変数（データベース接続）

```bash
export DATASOURCE_USERNAME=your_db_user
export DATASOURCE_PASSWORD=your_db_password
export DATASOURCE_HOST=localhost
export DATASOURCE_PORT=3306
export DATASOURCE_NAME=your_db_name
```

### Auth0認証用環境変数

```bash
# Auth0のドメイン（例: dev-example.us.auth0.com）
export AUTH0_DOMAIN=your-auth0-domain.auth0.com

# Auth0のAudience（例: https://your-api.example.com）
export AUTH0_AUDIENCE=your-api-audience
```

## 使用方法

### 1. アプリケーションの起動

```bash
# 依存関係の整理
go mod tidy

# ビルド
go build -o go-menu

# 実行（環境変数を設定して）
DATASOURCE_USERNAME=user \
DATASOURCE_PASSWORD=pass \
DATASOURCE_HOST=localhost \
DATASOURCE_PORT=3306 \
DATASOURCE_NAME=dbname \
AUTH0_DOMAIN=your-domain.auth0.com \
AUTH0_AUDIENCE=your-audience \
./go-menu
```

### 2. APIの利用

#### 認証不要なAPIの利用例

```bash
# ヘルスチェック
curl http://localhost:8080/v1/ping

# メニュー一覧取得
curl http://localhost:8080/v1/menus
```

#### 認証が必要なAPIの利用例

```bash
# お気に入り一覧取得
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:8080/v1/favorites

# お気に入り追加
curl -X POST \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"menu_id": 1}' \
     http://localhost:8080/v1/favorites

# お気に入り削除
curl -X DELETE \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:8080/v1/favorites/1
```

## エラーレスポンス

### 401 Unauthorized

以下の場合に返却されます：

- Authorizationヘッダーがない
- Bearerトークン形式が正しくない
- JWTトークンが無効
- トークンが期限切れ
- Audienceが一致しない
- Auth0 subが取得できない

```json
{
  "message": "Invalid token: token has invalid claims: token is expired"
}
```

### 403 Forbidden

- データベースにユーザーが存在しない場合

```json
{
  "message": "User not found"
}
```

### 500 Internal Server Error

- データベース接続エラーなどの内部エラー

```json
{
  "message": "Failed to get user: database connection error"
}
```

## セキュリティ対策

### 実装されているセキュリティ機能

1. **JWTトークン検証**: Auth0の公開鍵を使用してトークンの署名を検証
2. **有効期限チェック**: トークンの exp クレームを確認
3. **Audience検証**: 指定されたAudienceとトークンのAudienceを照合
4. **アルゴリズム検証**: RSA署名方式のみを許可
5. **個人識別情報の最小化**: Auth0 subのみを使用してユーザーを識別

### 注意事項

- Auth0ドメインとAudienceは正確に設定してください
- JWTトークンはHTTPS通信でのみ送信することを推奨
- データベース接続情報は環境変数で管理し、ソースコードにハードコードしないでください

## トラブルシューティング

### よくある問題

1. **"User not found" エラー**
   - データベースにAuth0 subに対応するユーザーレコードが存在しない
   - 最初にユーザー作成エンドポイントを呼び出してユーザーを作成してください

2. **"Invalid audience" エラー**
   - AUTH0_AUDIENCEとJWTトークンのaudienceクレームが一致しない
   - Auth0の設定を確認してください

3. **"Unable to find appropriate key" エラー**
   - Auth0のJWKS（JSON Web Key Set）からキーが見つからない
   - AUTH0_DOMAINが正しく設定されているか確認してください

## 実装詳細

### アーキテクチャ

```
middleware/auth.go           # Auth0認証ミドルウェア
├── JWTトークン抽出
├── Auth0公開鍵取得
├── トークン検証
└── ユーザー情報取得

handler/favoriteHandler.go  # お気に入り機能ハンドラー
├── お気に入り追加
├── お気に入り削除
└── お気に入り一覧取得

resource/user/userModel.go  # ユーザーデータ操作
├── Auth0 subでユーザー検索
├── お気に入り追加/削除
└── お気に入り一覧取得
```

### 依存性注入

DIコンテナを使用して各コンポーネントを接続：

```go
// di/wire.go
func InitFavoriteHandler() *handler.FavoriteHandler
func InitUserDriver() user.UserDriver
```

これにより、テスタビリティとモジュール性を確保しています。