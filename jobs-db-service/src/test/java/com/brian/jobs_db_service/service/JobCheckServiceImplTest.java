package com.brian.jobs_db_service.service;

import com.brian.jobs_db_service.dao.CompanyDao;
import com.brian.jobs_db_service.dao.JobDao;
import com.brian.jobs_db_service.model.entity.Job;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class JobCheckServiceImplTest {

    @Mock
    private JobDao jobDao;

    @Mock
    private CompanyDao companyDao;

    private JobCheckServiceImpl jobCheckService;

    @BeforeEach
    void setUp() {
        jobCheckService = new JobCheckServiceImpl(jobDao, companyDao);
    }

    @Test
    void jobAlreadyExists_returnsTrueWhenJobFound() {
        when(jobDao.getJob("job-1")).thenReturn(Optional.of(new Job("job-1", 1L)));

        assertThat(jobCheckService.jobAlreadyExists("job-1")).isTrue();
    }

    @Test
    void jobAlreadyExists_returnsFalseWhenJobMissing() {
        when(jobDao.getJob("job-1")).thenReturn(Optional.empty());

        assertThat(jobCheckService.jobAlreadyExists("job-1")).isFalse();
    }
}
