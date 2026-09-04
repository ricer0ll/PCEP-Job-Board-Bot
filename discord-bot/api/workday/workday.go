package workday

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

type WorkdayJobPostingRequest struct {
	AppliedFacets AppliedFacet `json:"appliedFacets"`
}

type WorkdayJobPostingResponse struct {
	Total       uint64 `json:"total"`
	JobPostings []WorkdayJobPosting
}
