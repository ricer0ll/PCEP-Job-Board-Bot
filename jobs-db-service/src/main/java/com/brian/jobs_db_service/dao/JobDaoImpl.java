package com.brian.jobs_db_service.dao;

import com.brian.jobs_db_service.model.entity.Company;
import com.brian.jobs_db_service.model.entity.Job;
import org.jdbi.v3.core.Jdbi;
import org.jdbi.v3.core.mapper.reflect.BeanMapper;
import org.jdbi.v3.core.statement.Query;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public class JobDaoImpl implements JobDao {

    private final Jdbi jdbi;

    public JobDaoImpl(Jdbi jdbi) {
        this.jdbi = jdbi;
    }

    @Override
    public Optional<Job> getJob(String jobId) {
        return jdbi.withHandle(handle -> {
            Query query = handle
                    .createQuery("select id, company_id from public.jobs where id = :job_id;")
                    .bind("job_id", jobId)
                    .registerRowMapper(BeanMapper.factory(Job.class));
            return query.mapTo(Job.class).findOne();
        });
    }

    @Override
    public Optional<Job> addJob(String jobId, Long companyId) {
        return jdbi.withHandle(handle -> {
            Query query = handle
                    .createQuery("select * from public.add_job(:job_id, :company_id);")
                    .bind("job_id", jobId)
                    .bind("company_id", companyId)
                    .registerRowMapper(BeanMapper.factory(Job.class));
            return query.mapTo(Job.class).findOne();
        });
    }
}
