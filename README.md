# VCNotify

自作のDiscord Voice Channel通知Botサンプル

## 概要
- Voice Channel参加／退出を検知し、指定時間帯のみ通知を行う

## 動作要件
- Go 1.21以上
- Discord Botトークン

## セットアップ
1. `.env.example` をコピーして `.env` にリネーム
2. 環境変数を設定 (`DISCORD_TOKEN`, `GUILD_ID`, `CHANNEL_ID`, `SKIP_START`, `SKIP_END`)
3. Botをローカルで起動:
   ```bash
   go run bot/main.go
4. Dockerで起動:
	docker build -t vcnotify .
	docker run --env-file .env vcnotify

