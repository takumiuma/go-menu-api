# プロジェクト構造とアーキテクチャ

## フォルダ構成
```
go-menu/
├── main.go                    # エントリーポイント
├── domain/                    # ドメインモデル
│   └── domain.go             # Menu構造体定義
├── handler/                   # HTTPハンドラー層
│   └── todoHandler.go        # MenuHandler実装
├── usecase/                   # ビジネスロジック層
│   └── port/                 # ポートインターフェース
│       └── port.go           # MenuPort定義
├── gateway/                   # ゲートウェイ層
│   └── gateway.go            # MenuGateway実装
├── resource/                  # データアクセス層
│   ├── dataAccess.go         # データベース接続
│   └── menu/                 # メニュー関連
├── di/                       # 依存性注入
│   └── wire.go               # Wire設定
└── router/                   # ルーティング設定
```

## アーキテクチャパターン

### クリーンアーキテクチャ
- **Domain**: ビジネスルールとエンティティ
- **Usecase**: アプリケーションロジック
- **Gateway**: 外部システムとの境界
- **Handler**: プレゼンテーション層
- **Resource**: データアクセス層

### 依存性の方向
```
Handler → Usecase → Port ← Gateway → Resource
```

## 命名規則

### 構造体・インターフェース
- 構造体: PascalCase (`MenuHandler`, `FavoriteDriver`)
- インターフェース: PascalCase + Interface suffix (`MenuPort`, `FavoriteDriver`)
- メソッド: PascalCase (`GetAll`, `CreateMenu`)

### ファイル・パッケージ
- ファイル: camelCase (`todoHandler.go`, `favoriteDriver.go`)
- パッケージ: lowercase (`domain`, `handler`, `usecase`)

### JSON タグ
- snake_case形式を使用 (`menu_id`, `menu_name`, `genre_ids`)

## 依存性注入パターン
- `Provide*` 関数でコンストラクタを提供
- Wireを使用した自動依存性注入
- インターフェースベースの疎結合設計

## エラーハンドリング
- 各層でエラーを適切に伝播
- HTTPレスポンスでは適切なステータスコードを返却
- エラーメッセージは `gin.H{"message": err.Error()}` 形式