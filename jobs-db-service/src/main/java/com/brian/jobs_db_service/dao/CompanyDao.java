package com.brian.jobs_db_service.dao;

import com.brian.jobs_db_service.model.entity.Company;

import java.math.BigInteger;

public interface CompanyDao {
    public Company getCompanyById(Long id);
    public Company addCompany(String name);
}
