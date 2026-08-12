package workday

import (
	"fmt"
	"testing"

	"github.com/ricer0ll/pcep-job-board/discord-bot/api/workday/dto"
)

func TestIsRelevantRole(t *testing.T) {
	client := &WorkdayClient{}

	cases := []struct {
		role     string
		expected bool
	}{
		{"Senior Software Engineer", true},
		{"Cloud Architect", true},
		{"Backend Developer", true},
		{"Product Owner", false},
		{"HR Manager", false},
		{"Manager", false},
	}

	for _, testCase := range cases {
		if client.isRelevantRole(testCase.role) != testCase.expected {
			t.Errorf("%s did not match the expected: %t", testCase.role, testCase.expected)
		}
	}
}

func TestGenerateNewJobPostingEmbedWorkday(t *testing.T) {
	client := &WorkdayClient{}

	companyName := "Contoso LLC"
	workdayURL := "https://workday.com"
	jobPosting := &dto.WorkdayJobPosting{
		Title:         "Software Engineer",
		ExternalPath:  "/test",
		LocationsText: "Portland, OR",
		PostedOn:      "Today",
		BulletFields:  []string{"ABC123"},
	}

	embed := client.generateNewJobPostingEmbed(jobPosting, companyName, workdayURL)

	exptectedTitle := fmt.Sprintf("New Job Posting from %s!", companyName)
	expectedDescription := fmt.Sprintf("Position: **%s**\nLocation: %s", jobPosting.Title, jobPosting.LocationsText)
	exptedURL := workdayURL + fmt.Sprintf("%s", jobPosting.ExternalPath)

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
