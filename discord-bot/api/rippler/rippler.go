package rippler

type RipplerJobPosting struct {
	JobTitle string `json:"job_title"`
	Location string `json:"location"`
}

type RipplerJobPostingRequest struct {
	URL string `json:"url"`
}

type RipplerJobPostingResponse struct {
	Jobs []RipplerJobPosting `json:"jobs"`
}
