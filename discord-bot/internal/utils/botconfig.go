package utils

import (
	"log/slog"
	"os"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
)

var (
	channelID string
	loaded    bool
)

func GetBotConfig() bot.ConfigOpt {
	configOpts := bot.WithGatewayConfigOpts(
		gateway.WithIntents(
			gateway.IntentGuilds,
			gateway.IntentGuildMessages,
		),
	)

	return configOpts
}

func GetDiscordChannelID() string {
	if loaded {
		return channelID
	}

	channelID = os.Getenv("CHANNEL_ID")
	if channelID == "" {
		slog.Error("Unable to load CHANNEL_ID environment variable.")
	}

	loaded = true
	return channelID
}
