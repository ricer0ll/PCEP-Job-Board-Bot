package rippler

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
	jobsCache       map[string][]dto.RipplerJobPosting = make(map[string][]dto.RipplerJobPosting)
	companyJsonPath string                             = filepath.Join("internal", "clients", "rippler", "companies.json")
)

const webscraperServiceUrl = "http://webscraper:8000/greenhouse/jobs"

type RipplerClient struct {
	restyClient *resty.Client
}

func NewRipplerClient(restyClient *resty.Client) *RipplerClient {
	return &RipplerClient{
		restyClient: restyClient,
	}
}

func (r RipplerClient) InitJobsCache() {
	companies, err := loadCompanies(companyJsonPath)
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("Loaded %d companies from Rippler config", len(companies)))

	for _, company := range companies {
		companyName := company.Name
		url := company.URL

		jobPostings, err := r.getRipplerJobPostings(url)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		jobsCache[companyName] = append(jobsCache[companyName], jobPostings.Jobs...)
	}
}

func (r RipplerClient) GetNewJobPostings(client *bot.Client) {
	companies, err := loadCompanies(companyJsonPath)
	if err != nil {
		panic("Unable to load companies")
	}

	for _, company := range companies {
		resp, err := r.getRipplerJobPostings(company.URL)
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
				r.notifyNewJob(client, &job, company.Name, company.URL) // notify on discord if new job
				jobsCache[company.Name] = append(jobsCache[company.Name], job)
			}
		}
	}

}

func (r RipplerClient) getRipplerJobPostings(url string) (*dto.RipplerJobPostingResponse, error) {
	request := dto.RipplerJobPostingRequest{
		URL: url,
	}
	resp := dto.RipplerJobPostingResponse{}

	result, err := r.restyClient.R().
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
		return nil, err
	}

	return &resp, nil
}

func (r RipplerClient) notifyNewJob(client *bot.Client, jobPosting *dto.RipplerJobPosting, company string, careerUrl string) {
	embed := r.generateNewJobPostingEmbed(jobPosting, company, careerUrl)
	client.Rest.CreateMessage(
		snowflake.MustParse(utils.GetDiscordChannelID()),
		discord.NewMessageCreate().WithEmbeds(embed),
	)

}

func (r RipplerClient) generateNewJobPostingEmbed(jobPosting *dto.RipplerJobPosting, company string, careerUrl string) discord.Embed {
	var title string = fmt.Sprintf("New Job Posting from %s!", company)
	var description string = fmt.Sprintf("Position: **%s**\nLocation: %s", jobPosting.JobTitle, jobPosting.Location)
	var url string = careerUrl

	embed := discord.NewEmbed().
		WithTitle(title).
		WithDescription(description).
		WithURL(url)

	return embed
}
