# API Gateway

## 用途
管理しているAPIへのリクエストをプロキシするゲートウェイです。

## 環境
以下の環境で動作を確認しています。
- Ubuntu 20.04.2
- go 1.16.4
- redis-server 6.2.3 (redis使用の場合)
- localstack 0.12.17 (dynamoDB使用の場合)

## 準備
以下の環境変数の設定が必要です。
- `READTIMEOUT`
    - リクエストのヘッダやボディを読む際のタイムアウト(sec)(ex. 5)
- `READHEADERTIMEOUT`
    - リクエストのヘッダを読む際のタイムアウト(sec)(ex. 5)
- `WRITETIMEOUT`
    - リクエストボディの読み込みからレスポンスの書き込みまでのタイムアウト(sec)(ex. 20)
- `IDLETIMEOUT`
    - リクエストを待つ際のタイムアウト(sec)(ex. 5)
- `MAXHEADERBYTES`
    - リクエストヘッダの容量の許可される最大値(ex. 1<<20)
- `LOG_PATH`
    - ログファイル(CSV形式)の出力先パス(ex. ./log.csv)
- `DB_TYPE`
    - 使用するDB。`REDIS`または`DYNAMO`。

以下は使用するDBに応じて設定

redisの場合
- `REDIS_HOST`
    - redisのホストアドレス(ex. localhost)
- `REDIS_PORT`
    - redisのlistenポート(ex. 6379)

dynamoDBの場合
- `DYNAMO_TABLE_API_FORWARDING`
    - APIルーティングを管理するテーブル名
- `DYNAMO_ENDPOINT`
    - localstack使用時に必要。接続先のホスト、ポート (ex. http://localhost:4566)

`source env.sh`でローカル実行用の環境変数を読み込むことが出来ます。

ログファイルは各列に日付(RFC3339形式)、APIキー、APIのパスをこの順で含んだCSV形式で作成されます。

## 実行

[Getting Started](../README_ja.md)を参照ください。
