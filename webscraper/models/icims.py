from pydantic import BaseModel

class IcimJob(BaseModel):
    job_title: str
    location: str

class IcimJobsRequest(BaseModel):
    url: str

class IcimJobsResponse(BaseModel):
    jobs: list[IcimJob]