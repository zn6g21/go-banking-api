# go-banking-api
[![CI](https://github.com/zn6g21/go-banking-api/actions/workflows/ci.yml/badge.svg)](https://github.com/zn6g21/go-banking-api/actions/workflows/ci.yml)

## 概要
Goで実装した銀行オープンAPI。OpenAPI 3.0 を起点に、アクセストークン再発行と口座情報取得を提供します。

## 背景・目的
私はメインフレームでの勘定系アプリのコーディング経験が3年半ありますが、
オープン系システムの開発には上流工程のみ携わっております。
そこで、バックエンドエンジニアのスキル習得および、転職活動のポートフォリオ作成を目的として、
ドメイン知識がある銀行オープンAPIをモダンな技術（Go/Docker/OpenAPI）で実装しました。

## 技術スタック
- Go 1.25.6 / Gin / GORM
- MySQL（Docker Compose）
- OpenAPI 3.0 + oapi-codegen
- GitHub Actions CI（lint / vulncheck / test / coverage）

## 主要機能
- Basic 認証の `/token` で refresh token を受け取り、access token を再発行
- Bearer 認証 + scope で `/accounts` を保護
- Health check: `GET /health`
- Swagger UI:
  - Docker Compose: http://localhost:8001/index.html
  - ローカル実行(APP_ENV=development): http://localhost:8080/swagger/index.html
- Docker Compose で API + MySQL + Swagger UI を一括起動

## アーキテクチャ
- entity / usecase / adapter / infrastructure に分割して責務を整理
    - 現職で担当するシステムは密結合のために変更影響が広く課題感を持っているため、影響範囲を局所化しやすいクリーンアーキテクチャを採用しました。

## API
Base URL: `http://localhost:8080/api/v1`

| Method | Path | Auth | Summary | Status |
| --- | --- | --- | --- | --- |
| GET | /accounts | Bearer | 口座情報取得 | ✅ |
| POST | /token | Basic | アクセストークン再発行 | ✅ |

## セットアップ（Docker Compose）
```sh
make docker-compose-up
```

### 停止
```sh
make docker-compose-down
```

## ローカル実行
```sh
make external-up
make run
```
APP_ENV=development のとき .env.development を読み込みます（必要なら作成）。

## OpenAPI / コード生成
api/openapi.yaml がAPI定義
`make generate-code-from-openapi` でコード生成

## テスト
```sh
make unittest
make test-cover
make integration-test
```

## CI
GitHub Actionsで lint / vulncheck / build / test / coverage を実施

## DBスキーマ
build/docker/external-apps/db/init.sql

## TODO
/transactions API 実装
