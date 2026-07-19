from pydantic import BaseModel
from .greenhouse import GreenhouseJob

class JamaJobsResponse(BaseModel):
    jobs: list[GreenhouseJob]

