package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	commands "github.com/ricer0ll/pcep-job-board/discord-bot/internal/bot/command"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/bot/event"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/utils"
)

func main() {
	token, err := utils.GetDiscordToken()
	if err != nil {
		panic(err)
	}

	configOpts := utils.GetBotConfig()

	client, err := disgo.New(
		*token,
		configOpts,
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds|cache.FlagMessages|cache.FlagMembers|cache.FlagRoles|cache.FlagChannels),
		),
		bot.WithEventListenerFunc(event.OnReady),
	)
	if err != nil {
		panic(err)
	}

	defer client.Close(context.TODO())

	if _, err = client.Rest.SetGlobalCommands(client.ApplicationID, commands.Commands); err != nil {
		slog.Error("error while registering commands", slog.Any("err", err))
	}

	// connect to the gateway
	if err = client.OpenGateway(context.TODO()); err != nil {
		panic(err)
	}

	// Start cron job to check jobs
	schedulerClient.InitCronJob(client)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}
