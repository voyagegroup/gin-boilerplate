# gin-boilerplate

[gin-gonic/gin](https://github.com/gin-gonic/gin) とReact(ES2015) を利用したTODOアプリです。TODOアプリのフロント実装は [tastejs/todomvc](https://github.com/tastejs/todomvc) を利用しています。

## 必要なもの

* Go 1.6以上
* MySQLクライアント
    * https://github.com/go-sql-driver/mysql のインストールに必要です。
* Node 6.x以上, npm 3.x以上
* (optional) Docker
    * MySQLをDockerコンテナで動かせるようにもしています。

## いろいろな操作

[Makefile]() をみてください。

```sh
# 環境設定をするには以下のようにします。
$ make deps

# テストの実行。ただしサーバサイドのみ。
$ make test

# インテグレーションテストの実行。テスト用のデータベースを利用したテストを実行します。
$ make integration-test

# eslintによるクライアント側アプリのlintを実行します。
$ make lint

# クライアントアプリの手動ビルド。
$ make build

# クライアント側コードを変更した場合に自動的にビルドします。`client` 以下のコードを保存すると自動的にコンパイルが走ります。
$ make watch
```

## データベースマイグレーション

[dbconfig.yml]() に置かれている設定によって接続するデータベースを指定しています。ここではデータベースマイグレーションの手順について説明します。

データベースマイグレーションでは [migrations]() 以下に置かれているSQLを実行するようにしています。

### (ローカル開発環境) 初回のみ

開発環境であれば、マイグレーションをする前にデータベースの作成が必要です。以下のようにするとデータベースを作成することができます。デフォルトではデータベース `treasure` を作成します。

    make migrate/init
    # インテグレーションテスト用のDBを作成するには以下のようにします
    make migrate/init DBNAME=test-treasure

パスワードをプロンプトで聞かれます。適宜パスワードを答えてください。後述するDocker環境であればデフォルトパスワードは `password` です。

### マイグレーションの実行

次にマイグレーションを実行します。

    # マイグレーションを実際には実行しないが適用されるであろうマイグレーションをプレビューすることができます
    make migrate/dry
    # マイグレーションを実行します
    make migrate/up
    # テストデータベースについて実行する場合にはENVを指定します
    make migrate/up ENV=test

マイグレーションを実行したあとに現在のデータベースの状態を確認するには以下のようにします。

    $ make migrate/status
    sql-migrate status
    +------------+-------------------------------+
    | MIGRATION  |            APPLIED            |
    +------------+-------------------------------+
    | 1_init.sql | 2016-07-21 05:18:52 +0000 UTC |
    +------------+-------------------------------+

マイグレーションには [rubenv/sql-migrate](https://github.com/rubenv/sql-migrate) を利用しています。詳しくはそちらをみるとよいでしょう。

## (optional) Docker環境について

MySQLをDockerで用意するようにしています。開発用に使う際には便利でしょう。もちろんlocalhostないしVM内にMySQLを自分で立てても構いません。

各環境でのDockerのインストールについては https://docs.docker.com/engine/installation/ をみてください。

開発用Docker環境を立ち上げるには以下のようにします。まずイメージをbuildします。

    make docker/build

[Dockerfile]() の内容を変更しないのであればbuildは一度だけしておけばよいでしょう。あとは以下のようにしてcontainerをスタートできます。

    make docker/start

コンテナに関する設定は [docker-compose.yml]() にまとめています。

## Acknowledgement

このアプリは学生エンジニア向けインターンシップ [Treasure](https://voyagegroup.com/internship/treasure/) のために用意されたサンプルアプリケーションです。
