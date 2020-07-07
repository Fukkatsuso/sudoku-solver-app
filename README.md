# Sudoku Solver App
任意の数独問題を解いてくれるWebアプリ


## Start
```sh
$ docker-compose up
```
の後，[http://localhost:8080]を開く


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
  - GCP_SA_EMAIL: サービスアカウントのemail
  - GCP_SA_KEY: サービスアカウントのJSON鍵をBase64エンコードする
    - (mac) `$ openssl base64 -in {key file} | pbcopy`
- 指定したGitHubのブランチにpush
- GCP > Cloud Run > サービス名 > 権限 > 追加 > 設定(allUsers, Cloud Run 起動元)


## TODO
- pushからデプロイ，マージまでのフローを整理
- viewsをコンポーネント分割