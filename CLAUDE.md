# shop

注文（order）と在庫（inventory）の2つの境界づけられたコンテキストからなるサンプルEC。

## アーキテクチャ

- `internal/<context>/domain`: 集約・値オブジェクト・ドメインエラー
- `internal/<context>/app`: アプリケーションサービス（ユースケース）
- `internal/<context>/infra`: リポジトリ実装

## ルール

- コンテキスト間の連携は、相手コンテキストの `app` パッケージが公開するサービス経由で行う
- 集約の状態変更は集約のメソッドを通して行う
- 金額は円単位の int64 で扱う
