package main

import (
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/greenhouse"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/jobsdb"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/rippler"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/scheduler"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/workday"
	"resty.dev/v3"
)

var (
	restyClient      *resty.Client                = resty.New()
	jobsDbClient     *jobsdb.JobsDbClient         = jobsdb.NewJobsDbClient(restyClient)
	workdayClient    *workday.WorkdayClient       = workday.NewWorkdayClient(restyClient, jobsDbClient)
	greenhouseClient *greenhouse.GreenhouseClient = greenhouse.NewGreenhouseClient(restyClient, jobsDbClient)
	ripplerClient    *rippler.RipplerClient       = rippler.NewRipplerClient(restyClient, jobsDbClient)
	schedulerClient  *scheduler.SchedulerClient   = scheduler.NewSchedulerClient(workdayClient, greenhouseClient, ripplerClient)
)
