package icims

type IcimsJobPosting struct {
	JobTitle string `json:"job_title"`
	Location string `json:"location"`
}

type IcimsJobPostingRequest struct {
	URL string `json:"url"`
}

type IcimsJobPostingResponse struct {
	Jobs []IcimsJobPosting `json:"jobs"`
}
