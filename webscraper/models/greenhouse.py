from pydantic import BaseModel

class GreenhouseJob(BaseModel):
    job_title: str
    location: str

class GreenhouseJobsRequest(BaseModel):
    url: str

class GreenhouseJobsResponse(BaseModel):
    jobs: list[GreenhouseJob]