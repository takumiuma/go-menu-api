# GraphQL API ドキュメント

このプロジェクトは、既存のREST APIに加えて、GraphQL APIも提供しています。

## エンドポイント

GraphQL APIは以下のエンドポイントで利用可能です：

```
POST http://localhost:8080/v1/graphql
GET  http://localhost:8080/v1/graphql (GraphQL Playground)
```

## スキーマ

### 型定義

#### Menu
```graphql
type Menu {
  menuId: Int!
  menuName: String!
  genreIds: [Int!]!
  categoryIds: [Int!]!
}
```

### クエリ (Query)

#### メニュー一覧の取得

```graphql
query {
  menus {
    menuId
    menuName
    genreIds
    categoryIds
  }
}
```

### ミューテーション (Mutation)

#### メニューの作成

```graphql
mutation {
  createMenu(input: {
    menuName: "新しいメニュー"
    genreIds: [1, 2]
    categoryIds: [3, 4]
  }) {
    menuId
    menuName
    genreIds
    categoryIds
  }
}
```

#### メニューの更新

```graphql
mutation {
  updateMenu(input: {
    menuId: 1
    menuName: "更新されたメニュー"
    genreIds: [1, 2, 3]
    categoryIds: [4, 5]
  }) {
    menuId
    menuName
    genreIds
    categoryIds
  }
}
```

#### メニューの削除

```graphql
mutation {
  deleteMenu(menuId: 1)
}
```

#### ジャンル関連の更新

```graphql
mutation {
  updateGenreRelations(input: {
    menuId: 1
    genreIds: [1, 2, 3]
  }) {
    menuId
    menuName
    genreIds
    categoryIds
  }
}
```

#### カテゴリ関連の更新

```graphql
mutation {
  updateCategoryRelations(input: {
    menuId: 1
    categoryIds: [4, 5, 6]
  }) {
    menuId
    menuName
    genreIds
    categoryIds
  }
}
```

## 使用例

### cURLでのリクエスト

```bash
curl -X POST http://localhost:8080/v1/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"{ menus { menuId menuName genreIds categoryIds } }"}'
```

### GraphQL Playgroundの使用

ブラウザで以下のURLにアクセスすると、GraphQL Playgroundが開きます：

```
http://localhost:8080/v1/graphql
```

Playgroundでは、スキーマの確認やクエリの実行をインタラクティブに行えます。

## アーキテクチャ

GraphQL実装は既存のREST API実装と同じレイヤーアーキテクチャを使用しています：

- **handler/graphqlHandler.go**: GraphQLハンドラー（スキーマ定義とリゾルバー）
- **graphql/resolver.go**: リゾルバーロジック（参考用）
- **usecase**: 既存のビジネスロジックを再利用
- **gateway**: 既存のゲートウェイ層を再利用
- **resource**: 既存のデータアクセス層を再利用

## 既存のREST APIとの比較

| 機能 | REST API | GraphQL API |
|------|----------|-------------|
| メニュー一覧取得 | GET /v1/menus | query { menus { ... } } |
| メニュー作成 | POST /v1/menus | mutation { createMenu(...) { ... } } |
| メニュー更新 | PUT /v1/menus/:id | mutation { updateMenu(...) { ... } } |
| メニュー削除 | DELETE /v1/menus/:id | mutation { deleteMenu(menuId: ...) } |
| ジャンル更新 | PATCH /v1/menus/:id/genres | mutation { updateGenreRelations(...) { ... } } |
| カテゴリ更新 | PATCH /v1/menus/:id/categories | mutation { updateCategoryRelations(...) { ... } } |

## 利点

GraphQL APIを使用することで、以下の利点があります：

1. **柔軟なデータ取得**: 必要なフィールドのみをリクエストできる
2. **単一エンドポイント**: すべての操作が /v1/graphql で完結
3. **型安全性**: スキーマによる厳密な型定義
4. **ドキュメント自動生成**: Playgroundでスキーマを確認可能
5. **バージョニング不要**: スキーマの進化をサポート

## 注意事項

- 現在、認証機能は実装されていません（メニューAPIのみ）
- お気に入り機能のGraphQL実装は今後追加予定
- エラーハンドリングは既存のusecaseレイヤーを利用
