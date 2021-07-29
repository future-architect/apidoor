# API Gateway

## 用途
管理しているAPIへのリクエストをプロキシするゲートウェイです。

## 環境
以下の環境で動作を確認しています。
- Ubuntu 20.04.2
- go 1.16.4
- redis-server 6.2.3

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
- `REDIS_HOST`
    - redisのホストアドレス(ex. localhost:6379)
- `LOG_PATH`
    - ログファイル(CSV形式)の出力先パス(ex. ./log.csv)

ログファイルは各列に日付(RFC3339形式)、APIキー、APIのパスをこの順で含んだCSV形式で作成されます。

## 実行
redis-serverの起動後、このREADMEのあるディレクトリで以下のコマンドを実行してください。
```
cd cmd/main
go run main.go
```