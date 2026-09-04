package icims

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"github.com/ricer0ll/pcep-job-board/discord-bot/api/icims"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/jobsdb"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/utils"
	"resty.dev/v3"
)

var (
	companyJsonPath string = filepath.Join("internal", "clients", "Icims", "companies.json")
)

const webscraperServiceUrl = "http://webscraper:8000/icims/jobs"

type DiscordRestClient interface {
	CreateMessage(channelID snowflake.ID, messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*discord.Message, error)
}

type IcimsClient struct {
	restyClient  *resty.Client
	jobsDbClient *jobsdb.JobsDbClient
}

func NewIcimsClient(restyClient *resty.Client, jobsDbClient *jobsdb.JobsDbClient) *IcimsClient {
	return &IcimsClient{
		restyClient:  restyClient,
		jobsDbClient: jobsDbClient,
	}
}

func (i IcimsClient) InitJobsCache(client *bot.Client) {
	companies, err := loadCompanies(companyJsonPath)
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("Loaded %d companies from Icims config", len(companies)))

	for _, company := range companies {
		companyName := company.Name
		url := company.URL

		jobPostings, err := i.getIcimsJobPostings(url)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		for _, job := range jobPostings.Jobs {
			exists, err := i.jobsDbClient.JobAlreadyExists(job.JobTitle, company.Name)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			if !exists {
				i.notifyNewJob(client.Rest, &job, companyName, url)
				i.jobsDbClient.AddJob(job.JobTitle, company.Name)
			}
		}
	}
}

func (i IcimsClient) GetNewJobPostings(client *bot.Client) {
	companies, err := loadCompanies(companyJsonPath)
	if err != nil {
		panic("Unable to load companies")
	}

	slog.Info("Getting Icims jobs")

	for _, company := range companies {
		resp, err := i.getIcimsJobPostings(company.URL)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		liveJobsPosting := resp.Jobs

		// check if job already exists. if not, add it to db and notify
		for _, job := range liveJobsPosting {
			exists, err := i.jobsDbClient.JobAlreadyExists(job.JobTitle, company.Name)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			if !exists {
				i.notifyNewJob(client.Rest, &job, company.Name, company.URL)
				i.jobsDbClient.AddJob(job.JobTitle, company.Name)
			}
		}
	}
}

func (i IcimsClient) getIcimsJobPostings(url string) (*icims.IcimsJobPostingResponse, error) {
	request := icims.IcimsJobPostingRequest{
		URL: url,
	}
	resp := icims.IcimsJobPostingResponse{}

	result, err := i.restyClient.R().
		SetContentType("application/json").
		SetBody(request).
		SetResult(&resp).
		SetTimeout(1 * time.Minute).
		Post(webscraperServiceUrl)

	if err != nil {
		slog.Error("Failed to post request to webscraper")
		return nil, err
	}
	if result.IsStatusFailure() {
		slog.Error("Post request to webscraper status code not 200's")
		return nil, fmt.Errorf("webscraper request failed with status code: %d", result.StatusCode())
	}

	return &resp, nil
}

func (i IcimsClient) notifyNewJob(client DiscordRestClient, jobPosting *icims.IcimsJobPosting, company string, careerUrl string) {
	embed := i.generateNewJobPostingEmbed(jobPosting, company, careerUrl)
	client.CreateMessage(
		snowflake.MustParse(utils.GetDiscordChannelID()),
		discord.NewMessageCreate().WithEmbeds(embed),
	)

}

func (i IcimsClient) generateNewJobPostingEmbed(jobPosting *icims.IcimsJobPosting, company string, careerUrl string) discord.Embed {
	var title string = fmt.Sprintf("New Job Posting from %s!", company)
	var description string = fmt.Sprintf("Position: **%s**\nLocation: %s", jobPosting.JobTitle, jobPosting.Location)
	var url string = careerUrl

	embed := discord.NewEmbed().
		WithTitle(title).
		WithDescription(description).
		WithURL(url)

	return embed
}
