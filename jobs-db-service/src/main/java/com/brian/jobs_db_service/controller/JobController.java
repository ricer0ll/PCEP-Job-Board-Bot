package com.brian.jobs_db_service.controller;

import com.brian.jobs_db_service.dao.JobDao;
import com.brian.jobs_db_service.model.dto.job.AddJobRequest;
import com.brian.jobs_db_service.model.dto.job.AddJobResponse;
import com.brian.jobs_db_service.model.dto.job.JobExistsRequest;
import com.brian.jobs_db_service.model.dto.job.JobExistsResponse;
import com.brian.jobs_db_service.model.entity.Job;
import com.brian.jobs_db_service.service.AddJobService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

import java.util.Optional;

@RestController
@RequestMapping("/v1/job")
public class JobController {
    private final AddJobService addJobService;

    public JobController(AddJobService addJobService) {
        this.addJobService = addJobService;
    }

    @GetMapping
    public JobExistsResponse checkIfJobExists(@RequestBody JobExistsRequest request) {
        return new JobExistsResponse(true, "The Standard");
    }

    @PostMapping
    public ResponseEntity<AddJobResponse> addJobResponse(@RequestBody AddJobRequest request) throws Exception {
        Job addedJob = addJobService
                .addJob(request.getJobName(), request.getCompanyName())
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.CONFLICT, "Job already exists"));

        return ResponseEntity.ok(new AddJobResponse(addedJob.getId(), addedJob.getCompanyId()));
    }
}
