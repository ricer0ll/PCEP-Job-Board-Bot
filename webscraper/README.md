# Webscraper

Some companies might not use Workday. Therefore, it is not as easy as making a HTTP requests to the workday API to get the list of jobs. Futhermore, some websites are dynamic, using JavaScript/TypeScript code to dynamically inject HTML elements.  

So we will be using Playwright to simulate a browser so that we can fetch their job listings as if we were visiting their website.