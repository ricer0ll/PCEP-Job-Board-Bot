package workday

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/ricer0ll/pcep-job-board/discord-bot/api/workday/dto"
	"github.com/ricer0ll/pcep-job-board/discord-bot/internal/utils"
	"resty.dev/v3"
)

var (
	jobsCache       map[string][]dto.WorkdayJobPosting = make(map[string][]dto.WorkdayJobPosting)
	companyJsonPath string                             = filepath.Join("internal", "clients", "workday", "companies.json")
	relevantRoles   []string                           = []string{"developer", "engineer", "software", "architect", "cloud"}
)

type WorkdayClient struct {
	restyClient *resty.Client
}

func NewWorkdayClient(restyClient *resty.Client) *WorkdayClient {
	return &WorkdayClient{
		restyClient: restyClient,
	}
}

func (w WorkdayClient) InitJobsCache() {
	companies, err := loadCompanies(companyJsonPath)
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("Loaded %d companies from Workday config", len(companies)))

	jobsCache := make(map[string][]dto.WorkdayJobPosting)

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

		jobsCache[company.Name] = append(jobsCache[company.Name], jobs...)
	}
}

func (w WorkdayClient) GetNewJobPostings(client *bot.Client) {
	companies, err := loadCompanies(companyJsonPath)
	if err != nil {
		panic("Unable to load companies")
	}

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

		// add to id cache to check later (basically a set)
		cachedIDs := make(map[string]struct{}) // Bullet Fields = ID (sorta)
		for _, job := range jobsCache[company.Name] {
			cachedIDs[job.BulletFields[0]] = struct{}{}
		}

		for _, job := range liveJobs {
			_, ok := cachedIDs[job.BulletFields[0]]
			if !ok && w.isRelevantRole(job.Title) {
				w.notifyNewJob(client, &job, company.Name, company.WorkdayBaseURL) // notify on discord if new job
				jobsCache[company.Name] = append(jobsCache[company.Name], job)
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
) ([]dto.WorkdayJobPosting, error) {
	jobPostings := []dto.WorkdayJobPosting{}

	request := dto.WorkdayJobPostingRequest{
		AppliedFacets: dto.AppliedFacet{
			JobFamily:       jobFamily,
			JobFamilyGroup:  jobFamilyGroup,
			LocationCountry: locationCountry,
			Locations:       locations,
		},
	}

	resp := dto.WorkdayJobPostingResponse{}

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

func (w WorkdayClient) notifyNewJob(client *bot.Client, jobPosting *dto.WorkdayJobPosting, company string, workdayURL string) {
	embed := w.generateNewJobPostingEmbed(jobPosting, company, workdayURL)
	client.Rest.CreateMessage(
		snowflake.MustParse(utils.GetDiscordChannelID()),
		discord.NewMessageCreate().WithEmbeds(embed),
	)

}

func (w WorkdayClient) generateNewJobPostingEmbed(jobPosting *dto.WorkdayJobPosting, company string, workdayURL string) discord.Embed {
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
