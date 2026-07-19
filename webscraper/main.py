from fastapi import FastAPI, status
from playwright.async_api import async_playwright
from models.greenhouse import GreenhouseJob, GreenhouseJobsResponse, GreenhouseJobsRequest

app = FastAPI()

keywords = ["developer", "engineer", "software", "architect", "cloud"]

@app.get("/health", status_code=status.HTTP_200_OK)
async def health_check():
    return {"status": "healthy"}

@app.post("/greenhouse/jobs")
async def get_jama_jobs(request: GreenhouseJobsRequest) -> GreenhouseJobsResponse:
    async with async_playwright() as p:
        browser = await p.firefox.launch(headless=True)
        page = await browser.new_page()

        await page.goto(request.url)
        
        await page.wait_for_timeout(3000)

        resp = GreenhouseJobsResponse(jobs=[])

        # note: iframe elements are a pain in the ass wtf man
        target_frame = None
        for frame in page.frames:
            frame_count = await frame.locator("div.job-posts").count()
            if frame_count > 0:
                target_frame = frame
                break

        if target_frame:
            titles_locator = target_frame.locator("tr.job-post p.body--medium")
            locations_locator = target_frame.locator("tr.job-post p.body__secondary")
            total_jobs = await titles_locator.count()
            
            for i in range(total_jobs):
                job_title = await titles_locator.nth(i).text_content()
                location = await locations_locator.nth(i).text_content()
                print(f"{i + 1}. {job_title.strip()} | {location.strip()}")

                if any(keyword in job_title.lower() for keyword in keywords):
                    resp.jobs.append(GreenhouseJob(job_title=job_title, location=location))
        else:
            print("Error: Could not locate the job board iframe container.")

        await browser.close()

        return resp

