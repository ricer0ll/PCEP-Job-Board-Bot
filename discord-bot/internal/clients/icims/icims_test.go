package icims

import (
	"fmt"
	"testing"

	"github.com/ricer0ll/pcep-job-board/discord-bot/api/icims"
)

func TestGenerateNewJobPostingEmbedGreenhouse(t *testing.T) {
	client := &IcimsClient{}

	companyName := "Contoso LLC"
	careersURL := "https://contoso.com"
	jobPosting := &icims.IcimsJobPosting{
		JobTitle: "Cloud Engineer",
		Location: "Remote",
	}

	embed := client.generateNewJobPostingEmbed(jobPosting, companyName, careersURL)

	exptectedTitle := fmt.Sprintf("New Job Posting from %s!", companyName)
	expectedDescription := fmt.Sprintf("Position: **%s**\nLocation: %s", jobPosting.JobTitle, jobPosting.Location)
	exptedURL := careersURL

	if embed.Title != exptectedTitle {
		t.Errorf("Title does not match. \nExpected: %s \tGot: %s", exptectedTitle, embed.Title)
	}

	if embed.Description != expectedDescription {
		t.Errorf("Description does not match. \nExpected: %s \tGot: %s", expectedDescription, embed.Description)
	}

	if embed.URL != exptedURL {
		t.Errorf("URL does not match. \nExpected: %s \tGot: %s", exptedURL, embed.URL)
	}
}
