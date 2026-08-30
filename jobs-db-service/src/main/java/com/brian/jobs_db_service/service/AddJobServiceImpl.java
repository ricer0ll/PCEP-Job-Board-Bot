package com.brian.jobs_db_service.service;

import com.brian.jobs_db_service.dao.CompanyDao;
import com.brian.jobs_db_service.dao.JobDao;
import com.brian.jobs_db_service.model.entity.Company;
import com.brian.jobs_db_service.model.entity.Job;
import com.brian.jobs_db_service.utils.Config;
import com.brian.jobs_db_service.utils.HasherUtil;
import com.google.common.hash.Hashing;
import org.springframework.stereotype.Service;

import java.nio.charset.StandardCharsets;
import java.util.Optional;

@Service
public class AddJobServiceImpl implements AddJobService {
    private final CompanyDao companyDao;
    private final JobDao jobDao;

    public AddJobServiceImpl(CompanyDao companyDao, JobDao jobDao) {
        this.companyDao = companyDao;
        this.jobDao = jobDao;
    }

    @Override
    public Optional<Job> addJob(String jobTitle, String companyName) throws Exception {
        Long companyId;

        Optional<Company> companyOpt = companyDao.getCompanyByName(companyName);
        if (companyOpt.isEmpty()) {
            companyId = companyDao
                    .addCompany(companyName)
                    .orElseThrow(() -> new Exception("Failed to create company"))
                    .getId();
        } else {
            companyId = companyOpt.get().getId();
        }

        String newJobId = HasherUtil.getJobId(jobTitle, companyName);
        return jobDao.addJob(newJobId, companyId);
    }
}
