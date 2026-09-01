package scheduler

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/go-co-op/gocron/v2"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/greenhouse"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/rippler"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/workday"
)

type jobClient interface {
	InitJobsCache()
	GetNewJobPostings(client *bot.Client)
}

type SchedulerClient struct {
	targetClients     []jobClient
	nonWorkdayClients []jobClient
}

func NewSchedulerClient(
	workdayClient *workday.WorkdayClient,
	greenhouseClient *greenhouse.GreenhouseClient,
	ripplerClient *rippler.RipplerClient,
) *SchedulerClient {
	return &SchedulerClient{
		targetClients: []jobClient{
			workdayClient,
			greenhouseClient,
			ripplerClient,
		},
	}
}

func (s *SchedulerClient) InitCronJob(client *bot.Client) (gocron.Scheduler, error) {
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		slog.Warn("Failed to load America/Los_Angeles timezone, falling back to UTC", "error", err)
		location = time.UTC
	}

	scheduler, err := gocron.NewScheduler(gocron.WithLocation(location))
	if err != nil {
		return nil, fmt.Errorf("failed to create cron scheduler: %w", err)
	}

	for _, clientObj := range s.targetClients {
		c := clientObj
		c.InitJobsCache()

		_, err = scheduler.NewJob(
			gocron.CronJob("0 6-18 * * 1-5", false),
			gocron.NewTask(c.GetNewJobPostings, client),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to schedule posting job: %w", err)
		}
	}

	scheduler.Start()
	slog.Info("Started cron jobs successfully")
	return scheduler, nil
}
