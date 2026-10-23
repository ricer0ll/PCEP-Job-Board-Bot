package com.brian.jobs_db_service.controller;

import com.brian.jobs_db_service.dao.JobDao;
import com.brian.jobs_db_service.model.entity.Job;
import com.brian.jobs_db_service.service.AddJobService;
import com.brian.jobs_db_service.service.JobCheckService;
import com.brian.jobs_db_service.utils.HasherUtil;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.webmvc.test.autoconfigure.WebMvcTest;
import org.springframework.http.MediaType;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;

import java.util.Optional;

import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(JobController.class)
class JobControllerTest {
    private static final String JOB_JSON = "{\"jobTitle\":\"Engineer\",\"companyName\":\"Acme\"}";
    private static final String JOB_ID = HasherUtil.getJobId("Engineer", "Acme");

    @Autowired
    private MockMvc mockMvc;

    @MockitoBean
    private AddJobService addJobService;

    @MockitoBean
    private JobCheckService jobCheckService;

    @MockitoBean
    private JobDao jobDao;

    @Test
    void checkIfJobExists_returnsTrueWhenJobExists() throws Exception {
        when(jobCheckService.jobAlreadyExists(JOB_ID)).thenReturn(true);

        mockMvc.perform(post("/v1/job/check")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(JOB_JSON))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.exists").value(true));
    }

    @Test
    void checkIfJobExists_returnsFalseWhenJobMissing() throws Exception {
        when(jobCheckService.jobAlreadyExists(JOB_ID)).thenReturn(false);

        mockMvc.perform(post("/v1/job/check")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(JOB_JSON))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.exists").value(false));
    }

    @Test
    void addJob_returnsAddedJob() throws Exception {
        when(addJobService.addJob("Engineer", "Acme")).thenReturn(Optional.of(new Job(JOB_ID, 3L)));

        mockMvc.perform(post("/v1/job")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(JOB_JSON))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.jobId").value(JOB_ID))
                .andExpect(jsonPath("$.companyId").value(3));
    }

    @Test
    void addJob_returnsConflictWhenJobExists() throws Exception {
        when(addJobService.addJob("Engineer", "Acme")).thenReturn(Optional.empty());

        mockMvc.perform(post("/v1/job")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(JOB_JSON))
                .andExpect(status().isConflict());
    }

    @Test
    void getJob_returnsJob() throws Exception {
        when(jobDao.getJob("abc")).thenReturn(Optional.of(new Job("abc", 3L)));

        mockMvc.perform(get("/v1/job/abc"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.jobId").value("abc"))
                .andExpect(jsonPath("$.companyId").value(3));
    }

    @Test
    void getJob_returnsNotFoundWhenMissing() throws Exception {
        when(jobDao.getJob("abc")).thenReturn(Optional.empty());

        mockMvc.perform(get("/v1/job/abc"))
                .andExpect(status().isNotFound());
    }
}
