package dto

type AppliedFacet struct {
	JobFamily       []string `json:"jobFamily,omitempty"`
	JobFamilyGroup  []string `json:"jobFamilyGroup,omitempty"`
	Locations       []string `json:"locations,omitempty"`
	LocationCountry []string `json:"locationCountry,omitempty"`
}

type WorkdayJobPosting struct {
	Title         string   `json:"title"`
	ExternalPath  string   `json:"externalPath"`
	LocationsText string   `json:"locationsText"`
	PostedOn      string   `json:"postedOn"`
	BulletFields  []string `json:"bulletFields"`
}

type GreenhouseJobPosting struct {
	JobTitle string `json:"job_title"`
	Location string `json:"location"`
}

type WorkdayJobPostingRequest struct {
	AppliedFacets AppliedFacet `json:"appliedFacets"`
}

type WorkdayJobPostingResponse struct {
	Total       uint64 `json:"total"`
	JobPostings []WorkdayJobPosting
}

type GreenhouseJobPostingRequest struct {
	URL string `json:"url"`
}

type GreenhouseJobPostingResponse struct {
	Jobs []GreenhouseJobPosting `json:"jobs"`
}
