package greenhouse

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"github.com/ricer0ll/pcep-job-board/discord-bot/api/greenhouse"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/jobsdb"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/utils"
	"resty.dev/v3"
)

var (
	companyJsonPath string = filepath.Join("internal", "clients", "greenhouse", "companies.json")
)

const webscraperServiceUrl = "http://webscraper:8000/greenhouse/jobs"

type DiscordRestClient interface {
	CreateMessage(channelID snowflake.ID, messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*discord.Message, error)
}

type GreenhouseClient struct {
	restyClient  *resty.Client
	jobsDbClient *jobsdb.JobsDbClient
}

func NewGreenhouseClient(restyClient *resty.Client, jobsDbClient *jobsdb.JobsDbClient) *GreenhouseClient {
	return &GreenhouseClient{
		restyClient:  restyClient,
		jobsDbClient: jobsDbClient,
	}
}

func (g GreenhouseClient) InitJobsCache(client *bot.Client) {
	companies, err := loadCompanies(companyJsonPath)
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("Loaded %d companies from Greenhouse config", len(companies)))

	for _, company := range companies {
		companyName := company.Name
		url := company.URL

		jobPostings, err := g.getGreenhouseJobPostings(url)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		for _, job := range jobPostings.Jobs {
			g.jobsDbClient.AddJob(job.JobTitle, companyName)
		}
	}
}

func (g GreenhouseClient) GetNewJobPostings(client *bot.Client) {
	companies, err := loadCompanies(companyJsonPath)
	if err != nil {
		panic("Unable to load companies")
	}

	slog.Info("Getting greenhouse jobs")

	for _, company := range companies {
		resp, err := g.getGreenhouseJobPostings(company.URL)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		liveJobsPosting := resp.Jobs

		// check if job already exists. if not, add it to db and notify
		for _, job := range liveJobsPosting {
			exists, err := g.jobsDbClient.JobAlreadyExists(job.JobTitle, company.Name)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			if !exists {
				g.notifyNewJob(client.Rest, &job, company.Name, company.URL)
				g.jobsDbClient.AddJob(job.JobTitle, company.Name)
			}
		}
	}
}

func (g GreenhouseClient) getGreenhouseJobPostings(url string) (*greenhouse.GreenhouseJobPostingResponse, error) {
	request := greenhouse.GreenhouseJobPostingRequest{
		URL: url,
	}
	resp := greenhouse.GreenhouseJobPostingResponse{}

	result, err := g.restyClient.R().
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

func (g GreenhouseClient) notifyNewJob(client DiscordRestClient, jobPosting *greenhouse.GreenhouseJobPosting, company string, careerUrl string) {
	embed := g.generateNewJobPostingEmbed(jobPosting, company, careerUrl)
	client.CreateMessage(
		snowflake.MustParse(utils.GetDiscordChannelID()),
		discord.NewMessageCreate().WithEmbeds(embed),
	)

}

func (g GreenhouseClient) generateNewJobPostingEmbed(jobPosting *greenhouse.GreenhouseJobPosting, company string, careerUrl string) discord.Embed {
	var title string = fmt.Sprintf("New Job Posting from %s!", company)
	var description string = fmt.Sprintf("Position: **%s**\nLocation: %s", jobPosting.JobTitle, jobPosting.Location)
	var url string = careerUrl

	embed := discord.NewEmbed().
		WithTitle(title).
		WithDescription(description).
		WithURL(url)

	return embed
}
