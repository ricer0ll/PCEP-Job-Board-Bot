package com.brian.jobs_db_service.service;

import com.brian.jobs_db_service.dao.CompanyDao;
import com.brian.jobs_db_service.dao.JobDao;
import com.brian.jobs_db_service.model.entity.Company;
import com.brian.jobs_db_service.model.entity.Job;
import com.brian.jobs_db_service.utils.HasherUtil;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyLong;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class AddJobServiceImplTest {
    private static final String TITLE = "Engineer";
    private static final String COMPANY = "Acme";
    private static final String JOB_ID = HasherUtil.getJobId(TITLE, COMPANY);

    @Mock
    private CompanyDao companyDao;

    @Mock
    private JobDao jobDao;

    private AddJobServiceImpl addJobService;

    @BeforeEach
    void setUp() {
        addJobService = new AddJobServiceImpl(companyDao, jobDao);
    }

    @Test
    void addJob_usesExistingCompany() throws Exception {
        Job job = new Job(JOB_ID, 7L);
        when(companyDao.getCompanyByName(COMPANY)).thenReturn(Optional.of(new Company(7L, COMPANY)));
        when(jobDao.addJob(JOB_ID, 7L)).thenReturn(Optional.of(job));

        assertThat(addJobService.addJob(TITLE, COMPANY)).containsSame(job);
        verify(companyDao, never()).addCompany(any());
    }

    @Test
    void addJob_createsCompanyWhenMissing() throws Exception {
        Job job = new Job(JOB_ID, 9L);
        when(companyDao.getCompanyByName(COMPANY)).thenReturn(Optional.empty());
        when(companyDao.addCompany(COMPANY)).thenReturn(Optional.of(new Company(9L, COMPANY)));
        when(jobDao.addJob(JOB_ID, 9L)).thenReturn(Optional.of(job));

        assertThat(addJobService.addJob(TITLE, COMPANY)).containsSame(job);
        verify(companyDao).addCompany(COMPANY);
    }

    @Test
    void addJob_throwsWhenCompanyCreationFails() {
        when(companyDao.getCompanyByName(COMPANY)).thenReturn(Optional.empty());
        when(companyDao.addCompany(COMPANY)).thenReturn(Optional.empty());

        assertThatThrownBy(() -> addJobService.addJob(TITLE, COMPANY))
                .hasMessage("Failed to create company");
        verify(jobDao, never()).addJob(any(), anyLong());
    }

    @Test
    void addJob_returnsEmptyWhenJobAlreadyExists() throws Exception {
        when(companyDao.getCompanyByName(COMPANY)).thenReturn(Optional.of(new Company(7L, COMPANY)));
        when(jobDao.addJob(JOB_ID, 7L)).thenReturn(Optional.empty());

        assertThat(addJobService.addJob(TITLE, COMPANY)).isEmpty();
    }
}
