from pydantic import BaseModel

class GreenhouseJob(BaseModel):
    title: str
    location: str