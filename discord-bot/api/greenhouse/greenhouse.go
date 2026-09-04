package greenhouse

type GreenhouseJobPosting struct {
	JobTitle string `json:"job_title"`
	Location string `json:"location"`
}

type GreenhouseJobPostingRequest struct {
	URL string `json:"url"`
}

type GreenhouseJobPostingResponse struct {
	Jobs []GreenhouseJobPosting `json:"jobs"`
}
