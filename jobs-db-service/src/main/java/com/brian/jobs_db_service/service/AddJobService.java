package com.brian.jobs_db_service.service;

import com.brian.jobs_db_service.model.entity.Job;

import java.util.Optional;

public interface AddJobService {
    public Optional<Job> addJob(String jobTitle, String companyName) throws Exception;
}
