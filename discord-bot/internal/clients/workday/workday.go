package workday

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"github.com/ricer0ll/pcep-job-board/discord-bot/api/workday"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/clients/jobsdb"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/utils"
	"resty.dev/v3"
)

var (
	companyJsonPath string   = filepath.Join("internal", "clients", "workday", "companies.json")
	relevantRoles   []string = []string{"developer", "engineer", "software", "architect", "cloud"}
)

type DiscordRestClient interface {
	CreateMessage(channelID snowflake.ID, messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*discord.Message, error)
}

type WorkdayClient struct {
	restyClient  *resty.Client
	jobsDbClient *jobsdb.JobsDbClient
}

func NewWorkdayClient(restyClient *resty.Client, jobsDbClient *jobsdb.JobsDbClient) *WorkdayClient {
	return &WorkdayClient{
		restyClient:  restyClient,
		jobsDbClient: jobsDbClient,
	}
}

func (w WorkdayClient) InitJobsCache() {
	companies, err := loadCompanies(companyJsonPath)
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("Loaded %d companies from Workday config", len(companies)))

	for _, company := range companies {
		jobs, err := w.getWorkdayJobPostings(
			company.WorkdayRequestURL,
			company.JobFamily,
			company.JobFamilyGroup,
			company.LocationCountry,
			company.Locations,
		)
		if err != nil {
			slog.Error(fmt.Sprintf("Unable to get job postings from %s", company.Name))
			continue
		}

		for _, job := range jobs {
			w.jobsDbClient.AddJob(job.Title, company.Name)
		}
	}
}

func (w WorkdayClient) GetNewJobPostings(client *bot.Client) {
	companies, err := loadCompanies(companyJsonPath)
	if err != nil {
		panic("Unable to load companies")
	}

	slog.Info("Getting workday jobs")

	for _, company := range companies {
		// get workday job postings
		liveJobs, err := w.getWorkdayJobPostings(
			company.WorkdayRequestURL,
			company.JobFamily,
			company.JobFamilyGroup,
			company.LocationCountry,
			company.Locations,
		)
		if err != nil {
			slog.Error(fmt.Sprintf("Unable to get job postings from %s", company.Name))
			continue
		}

		// check if job already exists. if not, add it to db and notify
		for _, job := range liveJobs {
			exists, err := w.jobsDbClient.JobAlreadyExists(job.Title, company.Name)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			if !exists {
				w.notifyNewJob(client.Rest, &job, company.Name, company.WorkdayBaseURL)
				w.jobsDbClient.AddJob(job.Title, company.Name)
			}
		}
	}
}

func (w WorkdayClient) getWorkdayJobPostings(
	url string,
	jobFamily []string,
	jobFamilyGroup []string,
	locationCountry []string,
	locations []string,
) ([]workday.WorkdayJobPosting, error) {
	jobPostings := []workday.WorkdayJobPosting{}

	request := workday.WorkdayJobPostingRequest{
		AppliedFacets: workday.AppliedFacet{
			JobFamily:       jobFamily,
			JobFamilyGroup:  jobFamilyGroup,
			LocationCountry: locationCountry,
			Locations:       locations,
		},
	}

	resp := workday.WorkdayJobPostingResponse{}

	result, err := w.restyClient.R().
		SetContentType("application/json").
		SetBody(request).
		SetResult(&resp).
		Post(url)

	if err != nil {
		slog.Error("Failed to post request to workday")
		return nil, err
	}
	if result.IsStatusFailure() {
		slog.Error("Post request to workday status code not 200's")
		return nil, fmt.Errorf("workday request failed with status code: %d", result.StatusCode())
	}

	jobPostings = resp.JobPostings
	return jobPostings, nil
}

func (w WorkdayClient) notifyNewJob(client DiscordRestClient, jobPosting *workday.WorkdayJobPosting, company string, workdayURL string) {
	embed := w.generateNewJobPostingEmbed(jobPosting, company, workdayURL)
	client.CreateMessage(
		snowflake.MustParse(utils.GetDiscordChannelID()),
		discord.NewMessageCreate().WithEmbeds(embed),
	)

}

func (w WorkdayClient) generateNewJobPostingEmbed(jobPosting *workday.WorkdayJobPosting, company string, workdayURL string) discord.Embed {
	var title string = fmt.Sprintf("New Job Posting from %s!", company)
	var description string = fmt.Sprintf("Position: **%s**\nLocation: %s", jobPosting.Title, jobPosting.LocationsText)
	var url string = workdayURL + fmt.Sprintf("%s", jobPosting.ExternalPath)

	embed := discord.NewEmbed().
		WithTitle(title).
		WithDescription(description).
		WithURL(url)

	return embed
}

func (w WorkdayClient) isRelevantRole(role string) bool {
	lowerCasedRole := strings.ToLower(role)

	for _, relevantRole := range relevantRoles {
		if strings.Contains(lowerCasedRole, relevantRole) {
			return true
		}
	}
	return false
}
