# API利用回数集計ロガー

## 用途
各APIをどのキーが何回利用したかを集計し、redisに書き込むアプリケーションです。

## 環境
以下の環境で動作を確認しています。
- Ubuntu 20.04.2
- go 1.16.4
- redis-server 6.2.3

## 準備
以下の環境変数の設定が必要です。
- `REDIS_HOST`
    - redisのホストアドレス(ex. localhost:6379)
- `LOG_PATH`
    - ログファイル(CSV形式)へのパス(ex. ./log.csv)

ログファイルは各列に日付(RFC3339形式)、APIキー、APIのパスをこの順で含んだCSV形式で作成してください。

## 実行
redis-serverを起動し、`PushLog()`を呼び出すとコードが実行されます。