package jobsdb

import (
	"fmt"
	"log/slog"

	"github.com/ricer0ll/pcep-job-board/discord-bot/api/jobsapi"
	"resty.dev/v3"
)

var (
	baseURL string = "http://jobs-db-service:8080/v1"
)

type JobsDbClient struct {
	restyClient *resty.Client
}

func NewJobsDbClient(restyClient *resty.Client) *JobsDbClient {
	return &JobsDbClient{
		restyClient: restyClient,
	}
}

func (j JobsDbClient) JobAlreadyExists(jobTitle string, companyName string) (bool, error) {
	req := jobsapi.JobExistsRequest{
		JobTitle:    jobTitle,
		CompanyName: companyName,
	}

	resp := jobsapi.JobExistsResponse{}

	result, err := j.restyClient.R().
		SetContentType("application/json").
		SetBody(req).
		SetResult(&resp).
		Post(baseURL + "/job/check")

	if err != nil {
		slog.Error("Failed to post request to check jobs")
		return false, err
	}
	if result.IsStatusFailure() {
		return false, fmt.Errorf("jobs db service request failed with status code: %d", result.StatusCode())
	}

	return resp.Exists, nil
}

func (j JobsDbClient) AddJob(jobTitle string, companyName string) {
	req := jobsapi.AddJobRequest{
		JobTitle:    jobTitle,
		CompanyName: companyName,
	}

	result, err := j.restyClient.R().
		SetContentType("application/json").
		SetBody(req).
		Post(baseURL + "/job")

	if err != nil {
		slog.Error("Failed to post request to jobs db service")
		return
	}
	if result.IsStatusFailure() {
		slog.Warn(fmt.Sprintf("%s from %s already exists", jobTitle, companyName))
		return
	}

	slog.Info("Added %s from %s", jobTitle, companyName)
}
