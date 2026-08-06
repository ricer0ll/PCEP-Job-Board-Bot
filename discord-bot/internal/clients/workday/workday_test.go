package workday

import "testing"

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
