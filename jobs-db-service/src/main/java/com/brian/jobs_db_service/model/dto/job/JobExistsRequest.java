package com.brian.jobs_db_service.model.dto.job;

import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter @Setter @NoArgsConstructor
public class JobExistsRequest {
    private String jobId;
    private String companyName;
}
