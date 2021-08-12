# 商材管理画面

## 用途
商材管理APIの操作をGUIで行うVue.jsアプリです。

## 環境
以下の環境で動作を確認しています。
- Ubuntu 20.04.2
- npm 6.14.4

## 準備
以下の環境変数を`.env`ファイル等に記述することが必要です。
- VUE_APP_API_BASE_URL
    - 用途: 商材管理APIのURL(ex. http://localhost:3000)

リポジトリのコードを変更せずローカルで実行する場合はexと同様に設定すると実行可能になります。

以下のコマンドを`package.json`があるディレクトリで実行し、セットアップを行ってください。
```
npm install
```

## 実行
各コマンドは`package.json`があるディレクトリで実行することが出来ます。データの取得のためには商材管理API(`management-api`ディレクトリ参照)の起動が必要ですが、後述のモックサーバを起動することによりこのディレクトリ内で完結することが出来ます。

### 開発用サーバの起動
```
npm run serve
```
コマンド実行後、http://localhost:8080 にアクセスすると画面を確認することが出来ます。

### プロジェクトのビルド
```
npm run build
```

### 単体テスト
```
npm run test:unit
```
[jest](https://jestjs.io/)を使用した単体テストを行っています。テストコードの詳細は`test/unit`ディレクトリをご覧ください。

### e2eテスト
```
npm run test:e2e
```
[cypress](https://www.cypress.io/)を用いたe2eテストを行っています。テストコードの詳細は`test/e2e`ディレクトリをご覧ください。

### lintの実行
```
npm run lint
```

### モックサーバの起動
```
npm run mock
```
[json-server](https://github.com/typicode/json-server)を用いたモックサーバです。返却値を変更する場合は`mock.json`を変更してください。
