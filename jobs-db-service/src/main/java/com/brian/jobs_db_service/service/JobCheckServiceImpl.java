package com.brian.jobs_db_service.service;

import com.brian.jobs_db_service.dao.CompanyDao;
import com.brian.jobs_db_service.dao.JobDao;
import com.brian.jobs_db_service.model.entity.Company;
import com.brian.jobs_db_service.model.entity.Job;
import org.springframework.stereotype.Service;

import java.util.Optional;

@Service
public class JobCheckServiceImpl implements JobCheckService {

    private final JobDao jobDao;
    private final CompanyDao companyDao;

    public JobCheckServiceImpl(JobDao jobDao, CompanyDao companyDao) {
        this.jobDao = jobDao;
        this.companyDao = companyDao;
    }

    @Override
    public Boolean jobAlreadyExists(String jobId, String companyName) {
        Optional<Job> jobOpt = jobDao.getJob(jobId);
        if (jobOpt.isEmpty()) {
            return false;
        }

        Optional<Company> companyOpt = companyDao.getCompanyByName(companyName);
        if (companyOpt.isEmpty()) {
            return false;
        }

        Job job = jobOpt.get();
        Company company = companyOpt.get();

        return job.getCompanyId().equals(company.getId());
    }
}
