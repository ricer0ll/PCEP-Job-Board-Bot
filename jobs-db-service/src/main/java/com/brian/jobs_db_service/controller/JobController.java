package com.brian.jobs_db_service.controller;

import com.brian.jobs_db_service.dao.JobDao;
import com.brian.jobs_db_service.model.dto.job.*;
import com.brian.jobs_db_service.model.entity.Job;
import com.brian.jobs_db_service.service.AddJobService;
import com.brian.jobs_db_service.service.JobCheckService;
import io.swagger.v3.oas.annotations.tags.Tag;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

import java.util.Optional;

@RestController
@RequestMapping("/v1/job")
@Tag(name = "Jobs")
public class JobController {
    private final AddJobService addJobService;
    private final JobCheckService jobCheckService;
    private final JobDao jobDao;

    public JobController(
            AddJobService addJobService,
            JobCheckService jobCheckService,
            JobDao jobDao
    ) {
        this.addJobService = addJobService;
        this.jobCheckService = jobCheckService;
        this.jobDao = jobDao;
    }

    @PostMapping("/check")
    public JobExistsResponse checkIfJobExists(@RequestBody JobExistsRequest request) {
        Boolean exists = jobCheckService.jobAlreadyExists(request.getJobId(), request.getCompanyName());

        return new JobExistsResponse(exists);
    }

    @PostMapping
    public ResponseEntity<AddJobResponse> addJobResponse(@RequestBody AddJobRequest request) throws Exception {
        Job addedJob = addJobService
                .addJob(request.getJobName(), request.getCompanyName())
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.CONFLICT, "Job already exists"));

        return ResponseEntity.ok(new AddJobResponse(addedJob.getId(), addedJob.getCompanyId()));
    }

    @GetMapping("/{job_id}")
    public ResponseEntity<GetJobResponse> getJobResponse(@PathVariable String job_id) {
        Job job = jobDao
                .getJob(job_id)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Job not found"));

        return ResponseEntity.ok(new GetJobResponse(job.getId(), job.getCompanyId()));
    }
}
