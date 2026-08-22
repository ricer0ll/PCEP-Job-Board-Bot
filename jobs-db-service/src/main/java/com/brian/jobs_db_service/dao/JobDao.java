package com.brian.jobs_db_service.dao;

import com.brian.jobs_db_service.model.entity.Job;

public interface JobDao {
    public Job getJob(String jobId);
    public void addJob(String jobId, Integer companyId);
}
