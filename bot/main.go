package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// 前回の VoiceState を管理するマップ
	lastVoice := make(map[string]string)

	// 環境変数取得
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN is not set")
	}
	skipStart, _ := strconv.Atoi(os.Getenv("SKIP_START"))
	skipEnd, _ := strconv.Atoi(os.Getenv("SKIP_END"))

	// Bot判定を無効化するか（デフォルト: false = 有効）
	disableBotFilter := os.Getenv("DISABLE_BOT_FILTER") == "true"

	// 除外するユーザーIDリスト
	excludeUserIDs := make(map[string]bool)
	if excludeIDs := os.Getenv("EXCLUDE_USER_IDS"); excludeIDs != "" {
		for _, id := range strings.Split(excludeIDs, ",") {
			id = strings.TrimSpace(id)
			if id != "" {
				excludeUserIDs[id] = true
			}
		}
	}

	// Discord セッション作成
	dg, err := discordgo.New(token)
	if err != nil {
		log.Fatal(err)
	}

	// VoiceStateUpdate イベントを受け取る Intent を有効化
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates

	dg.AddHandler(func(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
		// 前回の状態と比較して、同じチャンネルでの更新（ミュート／画面共有など）はスキップ
		prevChannel := lastVoice[vs.UserID]
		lastVoice[vs.UserID] = vs.ChannelID
		if prevChannel == vs.ChannelID && vs.ChannelID != "" {
			return
		}

		// 退出イベントは通知しない
		if vs.ChannelID == "" {
			return
		}

		// ユーザー情報を取得（Bot判定と表示名取得のため）
		user, err := s.User(vs.UserID)
		if err != nil {
			log.Printf("[WARN] failed to get user: %v\n", err)
			return
		}

		// Bot判定（環境変数で無効化可能）
		if !disableBotFilter && user.Bot {
			return
		}

		// 除外ユーザーIDリストのチェック
		if excludeUserIDs[vs.UserID] {
			return
		}

		// JST の時刻で判定
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		now := time.Now().In(jst)
		if skipStart <= now.Hour() && now.Hour() < skipEnd {
			return
		}

		// ユーザーのサーバーニックネーム（未設定時はユーザー名）
		member, err := s.GuildMember(vs.GuildID, vs.UserID)
		var displayName string
		if err == nil && member.Nick != "" {
			displayName = member.Nick
		} else {
			displayName = user.Username
		}

		// 通知送信
		channelID := os.Getenv("CHANNEL_ID")
		message := fmt.Sprintf("🔔 %s がボイスチャンネルに参加しました", displayName)
		if _, err := s.ChannelMessageSend(channelID, message); err != nil {
			log.Printf("[ERROR] failed to send message: %v\n", err)
		}
	})

	// WebSocket 接続をオープン
	if err := dg.Open(); err != nil {
		log.Fatal(err)
	}
	defer dg.Close()

	log.Println("VCNotify is running. Press CTRL-C to exit.")
	select {}
}
