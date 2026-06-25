# server

## 開発に必要なもの

- Go 1.26
- Docker
- [golangci-lint](https://golangci-lint.run/welcome/install/#local-installation) (静的解析)

### VSCode の場合

golangci-lint を VSCode で使うために、以下の手順が必要です。

- [{リポジトリルート}/.vscode/settings.template.json](../.vscode/settings.template.json) を `{リポジトリルート}/.vscode/settings.json` としてコピーする。

## ローカルで動かすだけの場合に必要なもの

- Docker

## コマンド

### 開発サーバーの起動

`.env` ファイルを作成して、以下の内容を記述する。

```dotenv
NS_MARIADB_DATABASE=h26s_02
NS_MARIADB_HOSTNAME=db
NS_MARIADB_PASSWORD=password
NS_MARIADB_PORT=3306
NS_MARIADB_USER=root
AWS_ACCESS_KEY_ID=rustfsadmin
AWS_SECRET_ACCESS_KEY=rustfsadmin
AWS_REGION=us-east-1
S3_BUCKET_NAME=h26s-02
S3_ENDPOINT=http://s3:9000
```

```bash
task up
```

ホットリロードの設定がされているので、コードを変更すると自動でコンテナのビルドが実行されて更新される。
ちょっと待つ必要がある。

- アプリ
- MariaDB
- Adminer
  - DB の中身を見たり操作したりできる
- rustfs
  - S3 互換のオブジェクトストレージ
  - http://localhost:9001 を開くと管理画面が見られる
  - アカウント、キー、ともに `rustfsadmin`

### 開発サーバーの停止

```bash
task down
```

### 静的解析

```bash
task lint
```
