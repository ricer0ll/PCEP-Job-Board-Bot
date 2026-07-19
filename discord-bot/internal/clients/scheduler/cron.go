package scheduler

import (
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/go-co-op/gocron/v2"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/greenhouse"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/workday"
)

type SchedulerClient struct {
	workdayClient    *workday.WorkdayClient
	greenhouseClient *greenhouse.GreenhouseClient
}

func NewSchedulerClient(
	workdayClient *workday.WorkdayClient,
	greenhouseClient *greenhouse.GreenhouseClient,
) *SchedulerClient {
	return &SchedulerClient{
		workdayClient:    workdayClient,
		greenhouseClient: greenhouseClient,
	}
}

func (s SchedulerClient) InitCronJob(client *bot.Client) gocron.Scheduler {
	location, _ := time.LoadLocation("America/Los_Angeles")
	scheduler, err := gocron.NewScheduler(gocron.WithLocation(location))
	if err != nil {
		panic("Failed to start cron scheduler!")
	}

	s.workdayClient.InitJobsCache()
	s.greenhouseClient.InitJobsCache()

	_, err = scheduler.NewJob(
		gocron.CronJob(
			"0 6-12 * * 1-5",
			false,
		),
		gocron.NewTask(s.workdayClient.GetNewJobPostings, client),
	)
	if err != nil {
		panic("Failed to create job for scheduler")
	}

	_, err = scheduler.NewJob(
		gocron.CronJob(
			"0 6-12 * * 1-5",
			false,
		),
		gocron.NewTask(s.greenhouseClient.GetNewJobPostings, client),
	)
	if err != nil {
		panic("Failed to create job for scheduler")
	}

	scheduler.Start()
	slog.Info("Started cron jobs")
	return scheduler
}
