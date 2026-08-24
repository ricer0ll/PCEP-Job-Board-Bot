package com.brian.jobs_db_service.model.dto.job;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter @Setter @NoArgsConstructor @AllArgsConstructor
public class GetJobIdRequest {
    private String jobTitle;
    private String companyName;
}
