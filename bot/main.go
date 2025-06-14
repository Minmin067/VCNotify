package main

import (
    "log"
    "os"
    "strconv"
    "time"

    "github.com/bwmarrin/discordgo"
)

func main() {
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
        // デバッグログ：イベント受信
        log.Printf("[DEBUG] VoiceStateUpdate UserID=%s ChannelID=%s\n", vs.UserID, vs.ChannelID)

        // 退出イベントは通知しない
        if vs.ChannelID == "" {
            log.Println("[DEBUG] leave event detected, skipping notification")
            return
        }

        // JST の時刻で判定
        jst := time.FixedZone("Asia/Tokyo", 9*60*60)
        now := time.Now().In(jst)
        log.Printf("[DEBUG] JST hour=%d (skip %d-%d)\n", now.Hour(), skipStart, skipEnd)
        if skipStart <= now.Hour() && now.Hour() < skipEnd {
            log.Println("[DEBUG] within skip window, skipping notification")
            return
        }
        log.Println("[DEBUG] outside skip window, sending notification")

        // 通知送信
        channelID := os.Getenv("CHANNEL_ID")
        message := "🔔 Voice channel activity detected"
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
