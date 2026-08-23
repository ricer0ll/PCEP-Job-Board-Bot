package com.brian.jobs_db_service.dao;

import com.brian.jobs_db_service.model.entity.Job;

import java.util.Optional;

public interface JobDao {
    public Optional<Job> getJob(String jobId);
    public Optional<Job> addJob(String jobId, Long companyId);
}
