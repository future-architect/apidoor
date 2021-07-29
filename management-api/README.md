# 商材管理API

## 用途
管理している商材(API)やそのセットに関する情報を提供し、また購入したものの利用状況を確認することが出来ます。

## 環境
以下の環境で動作を確認しています。
- Ubuntu 20.04.2
- go 1.16.4
- docker 20.10.6
- docker-compose 1.29.2

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

リポジトリのコードを変更せずローカルで実行する場合はexと同様に設定すると実行可能になります。

`docker-compose.yml`の`volumes`を、使用しているOSに関する記述以外コメントアウトしてください。

### Windowsユーザーのみ
dockerによるマウントがWSL上で出来ないため、`sql`ディレクトリをホストマシン内の任意の位置にコピーしてください。また、そのパスを`SQL_PATH`として環境変数に設定してください。

## 実行
このREADMEのあるディレクトリで以下のコマンドを実行してください。
```
docker-compose up -d
cd cmd/main
go run main.go
```