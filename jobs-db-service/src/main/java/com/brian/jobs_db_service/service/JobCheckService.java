package com.brian.jobs_db_service.service;

import org.springframework.stereotype.Service;

@Service
public interface JobCheckService {
    public Boolean jobAlreadyExists(String jobId, String companyName);
}
