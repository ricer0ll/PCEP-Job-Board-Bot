package jobsapi

type JobExistsRequest struct {
	JobTitle    string `json:"jobTitle"`
	CompanyName string `json:"companyName"`
}

type JobExistsResponse struct {
	Exists bool `json:"exists"`
}
