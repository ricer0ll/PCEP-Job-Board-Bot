package com.brian.jobs_db_service.model.entity;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.math.BigInteger;

@Getter @Setter @NoArgsConstructor @AllArgsConstructor
public class Company {
    private Long id;
    private String name;
}
