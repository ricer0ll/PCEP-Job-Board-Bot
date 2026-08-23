# Jobs Database Service
This service is mainly for any database operations. Mainly checking if jobs already exists or not. We do expose CRUD operations as well.  


## Endpoints
Here are the list of endpoints and their description.

| HTTP Method | Path | Request Body | Response |
|---|---|---|---|
| POST | `/v1/company` | `AddCompanyRequest` | `AddCompanyResponse` (200) / 409 Conflict if company exists |
| GET | `/v1/company/{company_id}` | – | `GetCompanyResponse` (200) / 404 Not Found |
| GET | `/v1/company/name/{company_name}` | – | `GetCompanyResponse` (200) / 404 Not Found |
| POST | `/v1/job/check` | `JobExistsRequest` | `JobExistsResponse` (200) |
| POST | `/v1/job` | `AddJobRequest` | `AddJobResponse` (200) / 409 Conflict if job exists |
| GET | `/v1/job/{job_id}` | – | `GetJobResponse` (200) / 404 Not Found |
| GET | `/root` | – | `GreetingResponse` (200) |
