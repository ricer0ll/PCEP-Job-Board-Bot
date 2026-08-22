package com.brian.jobs_db_service.model.entity;

import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.math.BigInteger;

@Getter @Setter @NoArgsConstructor
public class Job {
    private String id;
    private BigInteger companyId;
}
