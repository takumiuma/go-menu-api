# Copilot コーディングエージェント向けガイドライン

## プロジェクト概要

このプロジェクトは、メニューとそれに関連するジャンルやカテゴリを管理するGo言語ベースのRESTful APIサービスです。レストランや飲食店のメニューを管理し、メニューの作成・読み取り・更新・削除（CRUD）操作を提供します。

### 基本情報
- **言語**: Go 1.22.5以上
- **フレームワーク**: Gin Web Framework (v1.10.0)
- **データベース**: MySQL
- **ORM**: GORM (v1.25.7)
- **依存性注入**: Google Wire (v0.6.0)
- **プロジェクトサイズ**: 約900行のGoコード、12ファイル
- **実行ポート**: :8080 (デフォルト)

## 必須環境要件

### 依存関係
- Go 1.22.5以上（現在のバージョン: 1.24.7で動作確認済み）
- MySQLデータベース（動作中のMySQLサーバーが必要）

### 環境変数（必須）
アプリケーション実行前に以下の環境変数を設定する必要があります：
```bash
export DATASOURCE_USERNAME=your_db_user
export DATASOURCE_PASSWORD=your_db_password  
export DATASOURCE_HOST=localhost
export DATASOURCE_PORT=3306
export DATASOURCE_NAME=your_db_name
```

## ビルドと実行コマンド

### 依存関係の管理
```bash
# 依存関係のクリーンアップ（常に最初に実行）
go mod tidy

# 依存関係のダウンロード
go mod download
```

### コード検証
```bash
# コード品質チェック
go vet ./...

# フォーマット確認
gofmt -l .

# フォーマット修正（実行前に必ず確認）
gofmt -w .
```

### ビルド
```bash
# バイナリビルド（約15MBのバイナリが生成される）
go build -o go-menu

# または直接実行
go run main.go
```

### テスト
```bash
# テスト実行（現在テストファイルは存在しない）
go test ./...
```

**注意**: 現在プロジェクトにはテストファイルが存在しません。新しいテストを追加する場合は、既存のGoプロジェクトの慣例に従ってください。

### 実行
```bash
# 環境変数を設定してアプリケーション実行
DATASOURCE_USERNAME=user DATASOURCE_PASSWORD=pass DATASOURCE_HOST=localhost DATASOURCE_PORT=3306 DATASOURCE_NAME=dbname ./go-menu
```

**重要**: MySQLデータベースが利用できない場合、アプリケーションは起動時に失敗します（`connection refused`エラー）。

## プロジェクト構造とアーキテクチャ

### ディレクトリ構成
```
go-menu/
├── main.go                    # エントリーポイント
├── domain/                    # ドメインモデル層
│   └── domain.go             # Menu構造体定義
├── handler/                   # プレゼンテーション層
│   ├── todoHandler.go        # MenuHandlerの実装
│   └── system.go             # SystemHandler
├── usecase/                   # ビジネスロジック層
│   ├── usecase.go            # MenuUsecaseの実装
│   └── port/                 # ポートインターフェース
│       └── port.go           # MenuPort定義
├── gateway/                   # ゲートウェイ層
│   └── gateway.go            # MenuGateway実装
├── resource/                  # データアクセス層
│   ├── dataAccess.go         # データベース接続
│   └── menu/                 # メニュー関連
│       └── menuProvider.go   # MenuDriverの実装
├── di/                       # 依存性注入
│   ├── wire.go               # 手動DI設定
│   └── wireBuild.go          # Wire自動生成用
└── router/                   # ルーティング設定
    └── router.go             # Ginルーターの設定
```

### アーキテクチャパターン
クリーンアーキテクチャを採用しており、依存関係の方向は以下の通りです：
```
Handler → Usecase → Port ← Gateway → Resource
```

- **Domain**: ビジネスルールとエンティティ
- **Handler**: HTTPリクエスト/レスポンス処理
- **Usecase**: アプリケーションロジック
- **Gateway**: ドメインとインフラの境界
- **Resource**: データベースアクセス層

### API エンドポイント
```
GET    /v1/menus                           # メニュー一覧取得
POST   /v1/menus                           # メニュー作成
PUT    /v1/menus/:menu_id                  # メニュー更新
DELETE /v1/menus/:menu_id                  # メニュー削除
PATCH  /v1/menus/:menu_id/genres           # ジャンル関連更新
PATCH  /v1/menus/:menu_id/categories       # カテゴリ関連更新
```

## 重要な設定ファイル

### Go Modules (`go.mod`)
- プロジェクト名: `go-menu`
- 主要依存関係: Gin, GORM, Wire, MySQL driver, CORS

### VSCode設定 (`.vscode/launch.json`)
デバッグ設定が含まれています。Go拡張機能でのデバッグが可能です。

### ソースコード仕様

#### 命名規則
- **構造体/インターフェース**: PascalCase (`MenuHandler`, `MenuPort`)
- **メソッド**: PascalCase (`GetAll`, `CreateMenu`)
- **ファイル**: camelCase (`todoHandler.go`)
- **パッケージ**: lowercase (`domain`, `handler`)
- **JSON**: snake_case (`menu_id`, `menu_name`, `genre_ids`)

#### エラーハンドリング
- 各層でエラーを適切に伝播
- HTTPレスポンスでは適切なステータスコードを返却
- エラーメッセージ形式: `gin.H{"message": err.Error()}`

#### 依存性注入
- `Provide*` 関数でコンストラクタを提供
- Wireを使用した自動依存性注入（`di/wireBuild.go`）
- インターフェースベースの疎結合設計

## 開発時の注意事項

### コード変更時のワークフロー
1. 変更前に必ず `go mod tidy` を実行
2. コード変更後は `gofmt -w .` でフォーマット
3. `go vet ./...` で静的解析チェック
4. `go build -o go-menu` でビルド確認
5. 環境変数を設定してテスト実行

### 既知の問題と回避策
- フォーマットされていないコードが存在するため、変更前に `gofmt -w .` を実行することを推奨
- テストファイルが存在しないため、新機能追加時はテストも併せて作成することを検討
- MySQL接続が必須のため、開発環境では適切なデータベース設定が必要

### ファイル変更時の影響範囲
- **domain/domain.go**: 全層に影響（構造体変更時）
- **usecase/port/port.go**: Gateway層とUsecase層に影響
- **resource/dataAccess.go**: アプリケーション起動に直接影響
- **router/router.go**: エンドポイント変更時に影響
- **di/**: アプリケーション全体の依存関係に影響

## コーディングエージェント向け指示

### 信頼性向上のために
- このガイドラインの情報を信頼し、不完全または誤りがある場合のみ追加調査を実行してください
- ビルドコマンドは記載された順序で実行してください
- 環境変数の設定忘れによる実行失敗に注意してください
- コード変更後は必ずフォーマットとビルドチェックを実行してください

### 変更作業時の推奨手順
1. `go mod tidy && go mod download` で依存関係を整理
2. `gofmt -w .` でコードフォーマット
3. 変更を実装
4. `go vet ./...` で静的チェック
5. `go build -o go-menu` でビルド確認
6. 必要に応じて環境変数を設定して動作確認

このガイドラインに従うことで、ビルド失敗や検証エラーを最小限に抑え、効率的な開発が可能になります。