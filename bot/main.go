package main

import (
    "log"
    "os"
    "strconv"
    "time"

    "github.com/bwmarrin/discordgo"
)

func main() {
    token := os.Getenv("DISCORD_TOKEN")
    if token == "" {
        log.Fatal("DISCORD_TOKEN is not set")
    }
    skipStart, _ := strconv.Atoi(os.Getenv("SKIP_START"))
    skipEnd, _ := strconv.Atoi(os.Getenv("SKIP_END"))

    dg, err := discordgo.New(token)
    if err != nil {
        log.Fatal(err)
    }
    dg.AddHandler(func(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
        now := time.Now().Hour()
        if skipStart <= now && now < skipEnd {
            return  // 指定時間帯は通知しない
        }
        channelID := os.Getenv("CHANNEL_ID")
        message := "🔔 Voice channel activity detected"
        s.ChannelMessageSend(channelID, message)
    })

    if err := dg.Open(); err != nil {
        log.Fatal(err)
    }
    defer dg.Close()

    log.Println("VCNotify is running. Press CTRL-C to exit.")
    select {}
}