from pydantic import BaseModel

class RipplerJob(BaseModel):
    job_title: str
    location: str

class RipplerJobRequest(BaseModel):
    url: str

class RipplerJobResponse(BaseModel):
    jobs: list[RipplerJob]