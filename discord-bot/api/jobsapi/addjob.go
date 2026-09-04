package jobsapi

type AddJobRequest struct {
	JobTitle    string `json:"jobTitle"`
	CompanyName string `json:"companyName"`
}
