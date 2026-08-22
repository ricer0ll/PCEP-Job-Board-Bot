package com.brian.jobs_db_service.controller;

import com.brian.jobs_db_service.model.dto.job.JobExistsRequest;
import com.brian.jobs_db_service.model.dto.job.JobExistsResponse;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/v1/job")
public class JobController {
    @GetMapping
    public JobExistsResponse checkIfJobExists(@RequestBody JobExistsRequest request) {
        return new JobExistsResponse(true, "The Standard");
    }
}
