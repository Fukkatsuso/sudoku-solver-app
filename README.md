# Sudoku Solver App

[![Deploy](https://github.com/Fukkatsuso/sudoku-solver-app/actions/workflows/deploy.yml/badge.svg)](https://github.com/Fukkatsuso/sudoku-solver-app/actions/workflows/deploy.yml)

任意の数独問題を解いてくれるWebアプリ

## Start

```sh
$ docker-compose up
```

の後，<http://localhost:8080>を開く

## Deploy

### GCP Cloud Run

- Cloud Run APIの有効化
- [サービスアカウントを作成](https://cloud.google.com/iam/docs/creating-managing-service-accounts?hl=ja#iam-service-accounts-create-console)
  - Cloud Run 管理者
  - Cloud Storage 管理者
  - サービス アカウント ユーザ
- サービスアカウントのJSON鍵を生成
- リポジトリのSecrets
  - GCP_PROJECT: プロジェクトID
  - GCP_REGION: リージョン
  - GCP_SA_KEY: サービスアカウントのJSON鍵
- 指定したGitHubのブランチにpush

## TODO

- pushからデプロイ，マージまでのフローを整理
- viewsをコンポーネント分割
- 操作を戻す機能
- テーブル内の同じ数字を修飾する機能
- メモ書き機能
- OCRの精度向上
- dockerイメージの軽量化

## References

- [gocv/opencv - Docker container with OpenCV 4.x](https://hub.docker.com/r/gocv/opencv)
- [tesseract installation](https://github.com/tesseract-ocr/tessdoc/blob/main/Installation.md)
