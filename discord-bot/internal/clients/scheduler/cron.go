package scheduler

import (
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/go-co-op/gocron/v2"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/greenhouse"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/rippler"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/workday"
)

type SchedulerClient struct {
	workdayClient    *workday.WorkdayClient
	greenhouseClient *greenhouse.GreenhouseClient
	ripplerClient    *rippler.RipplerClient
}

func NewSchedulerClient(
	workdayClient *workday.WorkdayClient,
	greenhouseClient *greenhouse.GreenhouseClient,
	ripplerClient *rippler.RipplerClient,
) *SchedulerClient {
	return &SchedulerClient{
		workdayClient:    workdayClient,
		greenhouseClient: greenhouseClient,
		ripplerClient:    ripplerClient,
	}
}

func (s SchedulerClient) InitCronJob(client *bot.Client) gocron.Scheduler {
	location, _ := time.LoadLocation("America/Los_Angeles")
	scheduler, err := gocron.NewScheduler(gocron.WithLocation(location))
	if err != nil {
		panic("Failed to start cron scheduler!")
	}

	type jobClient interface {
		InitJobsCache()
		GetNewJobPostings(client *bot.Client)
	}

	targetClients := []jobClient{
		s.workdayClient,
		s.greenhouseClient,
		s.ripplerClient,
	}

	for _, c := range targetClients {
		c.InitJobsCache()

		_, err = scheduler.NewJob(
			gocron.CronJob("0 6-12 * * 1-5", false),
			gocron.NewTask(c.GetNewJobPostings, client),
		)
		if err != nil {
			panic("Failed to create job for scheduler")
		}
	}

	// Need to reset non workday clients as they dont use ids
	nonWorkdayClients := []jobClient{
		s.greenhouseClient,
		s.ripplerClient,
	}

	for _, c := range nonWorkdayClients {
		_, err = scheduler.NewJob(
			gocron.CronJob("0 0 * * 0", false),
			gocron.NewTask(c.InitJobsCache),
		)
	}

	scheduler.Start()
	slog.Info("Started cron jobs")
	return scheduler
}
