package greenhouse

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/ricer0ll/pcep-job-board/discord-bot/api/workday/dto"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/utils"
	"resty.dev/v3"
)

var (
	jobsCache       map[string][]dto.GreenhouseJobPosting = make(map[string][]dto.GreenhouseJobPosting)
	companyJsonPath string                                = filepath.Join("internal", "clients", "greenhouse", "companies.json")
)

const webscraperServiceUrl = "http://webscraper:8000/greenhouse/jobs"

type GreenhouseClient struct {
	restyClient *resty.Client
}

func NewGreenhouseClient(restyClient *resty.Client) *GreenhouseClient {
	return &GreenhouseClient{
		restyClient: restyClient,
	}
}

func (g GreenhouseClient) InitJobsCache() {
	companies, err := loadCompanies(companyJsonPath)
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("Loaded %d companies from Greenhouse config", len(companies)))

	jobsCache := make(map[string][]dto.GreenhouseJobPosting)

	for _, company := range companies {
		companyName := company.Name
		url := company.URL

		jobPostings, err := g.getGreenhouseJobPostings(url)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		jobsCache[companyName] = append(jobsCache[companyName], jobPostings.Jobs...)
	}
}

func (g GreenhouseClient) GetNewJobPostings(client *bot.Client) {
	companies, err := loadCompanies(companyJsonPath)
	if err != nil {
		panic("Unable to load companies")
	}

	for _, company := range companies {
		resp, err := g.getGreenhouseJobPostings(company.URL)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		liveJobsPosting := resp.Jobs

		// add job title to cache
		// (yes, ik this is bad, but greenhouse doesn't give me a way to get job's id...)
		cachedIDs := make(map[string]struct{})
		for _, job := range jobsCache[company.Name] {
			cachedIDs[job.JobTitle] = struct{}{}
		}

		for _, job := range liveJobsPosting {
			_, ok := cachedIDs[job.JobTitle]
			if !ok {
				g.notifyNewJob(client, &job, company.Name, company.URL) // notify on discord if new job
				jobsCache[company.Name] = append(jobsCache[company.Name], job)
			}
		}
	}
}

func (g GreenhouseClient) getGreenhouseJobPostings(url string) (*dto.GreenhouseJobPostingResponse, error) {
	request := dto.GreenhouseJobPostingRequest{
		URL: url,
	}
	resp := dto.GreenhouseJobPostingResponse{}

	result, err := g.restyClient.R().
		SetContentType("application/json").
		SetBody(request).
		SetResult(&resp).
		SetTimeout(1 * time.Minute).
		Post(webscraperServiceUrl)

	if err != nil {
		slog.Error("Failed to post request to workday")
		return nil, err
	}
	if result.IsStatusFailure() {
		slog.Error("Post request to workday status code not 200's")
		return nil, fmt.Errorf("webscraper request failed with status code: %d", result.StatusCode())
	}

	return &resp, nil
}

func (g GreenhouseClient) notifyNewJob(client *bot.Client, jobPosting *dto.GreenhouseJobPosting, company string, careerUrl string) {
	embed := g.generateNewJobPostingEmbed(jobPosting, company, careerUrl)
	client.Rest.CreateMessage(
		snowflake.MustParse(utils.GetDiscordChannelID()),
		discord.NewMessageCreate().WithEmbeds(embed),
	)

}

func (g GreenhouseClient) generateNewJobPostingEmbed(jobPosting *dto.GreenhouseJobPosting, company string, careerUrl string) discord.Embed {
	var title string = fmt.Sprintf("New Job Posting from %s!", company)
	var description string = fmt.Sprintf("Position: **%s**\nLocation: %s", jobPosting.JobTitle, jobPosting.Location)
	var url string = careerUrl

	embed := discord.NewEmbed().
		WithTitle(title).
		WithDescription(description).
		WithURL(url)

	return embed
}
