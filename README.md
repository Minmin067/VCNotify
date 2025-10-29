# VCNotify

自作のDiscord Voice Channel通知Botサンプル

## 概要
- Voice Channel参加を検知し、指定時間帯のみ通知を行う
- Botユーザーの参加通知を自動的に除外
- 特定ユーザーの参加通知を除外可能

## 動作要件
- Go 1.21以上
- Discord Botトークン

## セットアップ
1. `.env.example` をコピーして `.env` にリネーム
2. 環境変数を設定

### 必須環境変数
- `DISCORD_TOKEN`: Discord Botトークン
- `CHANNEL_ID`: 通知を送信するテキストチャンネルID
- `SKIP_START`: 通知をスキップする開始時刻（0-23の数字）
- `SKIP_END`: 通知をスキップする終了時刻（0-23の数字）

### オプション環境変数
- `DISABLE_BOT_FILTER`: Bot判定を無効化する場合に `true` を設定（デフォルト: 無効 = Botを自動除外）
- `EXCLUDE_USER_IDS`: 通知を除外するユーザーIDをカンマ区切りで指定（例: `123456789,987654321`）

3. Botをローカルで起動:
   ```bash
   go run bot/main.go
4. Dockerで起動:
	docker build -t vcnotify .
	docker run --env-file .env vcnotify

