# API Gateway

## 用途

Web APIへのリクエストをフォワードするゲートウェイサービスを構築するためのライブラリです。Goでは動的にデータソースを入れ替えることができないため、mainパッケージ内で静的にプラグインを切り替える思想で設計されています。

### cmd/localredisgateway

localhost内で起動することを目的に組み合わせられたエントリーポイントです。データソースにはRedisを利用します。

#### 環境変数の設定

* [ ] LOG_PATH
    - ログファイル(CSV形式)の出力先パス
    - デフォルト: ./log.csv
* [ ] REDIS_HOST
    - redisのホストアドレス
    - デフォルト: localhost
* [ ] REDIS_PORT
    - redisのlistenポート
    - デフォルト: 6379

### cmd/localdynamogateway

dynamoDBの場合
- `DYNAMO_TABLE_API_FORWARDING`
    - APIルーティングを管理するテーブル名
- `DYNAMO_ENDPOINT`
    - localstack使用時に必要。接続先のホスト、ポート (ex. http://localhost:4566)

`source env.sh`でローカル実行用の環境変数を読み込むことが出来ます。

ログファイルは各列に日付(RFC3339形式)、APIキー、APIのパスをこの順で含んだCSV形式で作成されます。

## 実行

[Getting Started](../README_ja.md)を参照ください。

## Data model

| name        | type   | memo | example                         |
|-------------|--------|------|---------------------------------|
| api_key     | string |      | key                             |
| path        | string |      | test                            |
| forward_url | string |      | http://test-server:3333/welcome |

