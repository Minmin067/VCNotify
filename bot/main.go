package main

import (
    "fmt"
    "log"
    "os"
    "strconv"
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

    // Discord セッション作成
    dg, err := discordgo.New(token)
    if err != nil {
        log.Fatal(err)
    }

    // VoiceStateUpdate イベントを受け取る Intent を有効化
    dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates

    dg.AddHandler(func(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
                // 前回の状態と比較して、同じチャンネルでの更新（ミュート／画面共有など）はスキップ
        prevChannel, _ := lastVoice[vs.UserID]
        lastVoice[vs.UserID] = vs.ChannelID
        if prevChannel == vs.ChannelID && vs.ChannelID != "" {
            return
        }

        // 退出イベントは通知しない
        if vs.ChannelID == "" {
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
            user, err := s.User(vs.UserID)
            if err == nil {
                displayName = user.Username
            } else {
                displayName = vs.UserID
            }
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
