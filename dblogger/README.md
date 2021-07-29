# API使用履歴データロガー

## 用途
APIの利用履歴(利用日時、利用のために使用したAPIキー、利用したAPIのパス)をデータベースに書き込むアプリケーションです。

## 環境
以下の環境で動作を確認しています。
- Ubuntu 20.04.2
- go 1.16.4

## 準備
以下の環境変数の設定が必要です。
- `DATABASE_DRIVER`
    - 用途: データベースのドライバ(ex. postgres)
- `DATABASE_HOST`
    - 用途: データベースのホストアドレス(ex. 127.0.0.1)
- `DATABASE_PORT`
    - 用途: データベースの接続ポート(ex. 5555)
- `DATABASE_USER`
    - 用途: データベースに接続するユーザ名(ex. root)
- `DATABASE_PASSWORD`
    - 用途: データベースに接続するパスワード(ex. password)
- `DATABASE_NAME`
    - 用途: データベース名(ex. root)
- `DATABASE_SSLMODE`
    - 用途: SSLを有効化するか(ex. disable)
- `LOG_PATH`
    - ログファイル(CSV形式)へのパス(ex. ./log.csv)

ログファイルは各列に日付(RFC3339形式)、APIキー、APIのパスをこの順で含んだCSV形式で作成してください。

## 実行
データベースを起動し、`PushLog()`を呼び出すとコードが実行されます。