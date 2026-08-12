from unittest.mock import AsyncMock, MagicMock, patch
import pytest
from httpx import ASGITransport, AsyncClient
from main import app


@pytest.fixture
def mock_playwright():
    with patch("main.async_playwright") as mock_pw:
        pw_instance = AsyncMock()

        cm = AsyncMock()
        cm.__aenter__.return_value = pw_instance
        mock_pw.return_value = cm

        browser = AsyncMock()
        pw_instance.firefox.launch = AsyncMock(return_value=browser)

        page = AsyncMock()
        page.wait_for_timeout = AsyncMock()
        browser.new_page = AsyncMock(return_value=page)

        yield page


@pytest.fixture
async def async_client():
    async with AsyncClient(
        transport=ASGITransport(app=app), base_url="http://test"
    ) as client:
        yield client


@pytest.mark.anyio
async def test_health_check(async_client):
    response = await async_client.get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "healthy"}


@pytest.mark.anyio
async def test_get_greenhouse_jobs_success(async_client, mock_playwright):
    mock_frame = AsyncMock()
    mock_playwright.frames = [mock_frame]

    job_posts_locator = AsyncMock()
    job_posts_locator.count = AsyncMock(return_value=1)

    titles_locator = AsyncMock()
    titles_locator.count = AsyncMock(return_value=2)

    title1_mock = AsyncMock()
    title1_mock.text_content = AsyncMock(return_value=" Software Engineer ")
    title2_mock = AsyncMock()
    title2_mock.text_content = AsyncMock(return_value=" Sales Representative ")

    titles_locator.nth = MagicMock(side_effect=[title1_mock, title2_mock])

    locations_locator = AsyncMock()
    locations_locator.count = AsyncMock(return_value=2)

    loc1_mock = AsyncMock()
    loc1_mock.text_content = AsyncMock(return_value=" Remote ")
    loc2_mock = AsyncMock()
    loc2_mock.text_content = AsyncMock(return_value=" New York ")

    locations_locator.nth = MagicMock(side_effect=[loc1_mock, loc2_mock])

    def locator_side_effect(selector):
        if selector == "div.job-posts":
            return job_posts_locator
        if selector == "tr.job-post p.body--medium":
            return titles_locator
        if selector == "tr.job-post p.body__secondary":
            return locations_locator
        return AsyncMock()

    mock_frame.locator = MagicMock(side_effect=locator_side_effect)

    payload = {"url": "https://boards.greenhouse.io/testcompany"}
    response = await async_client.post("/greenhouse/jobs", json=payload)

    assert response.status_code == 200
    data = response.json()

    assert len(data["jobs"]) == 1
    assert data["jobs"][0]["job_title"] == "Software Engineer"
    assert data["jobs"][0]["location"] == "Remote"


@pytest.mark.anyio
async def test_get_rippling_jobs_frame_not_found(async_client, mock_playwright):
    mock_frame = AsyncMock()

    text_locator = AsyncMock()
    text_locator.count = AsyncMock(return_value=0)
    mock_frame.get_by_text = MagicMock(return_value=text_locator)

    mock_playwright.frames = [mock_frame]

    payload = {"url": "https://ats.rippling.com/testcompany"}
    response = await async_client.post("/rippling/jobs", json=payload)

    assert response.status_code == 200
    assert response.json() == {"jobs": []}