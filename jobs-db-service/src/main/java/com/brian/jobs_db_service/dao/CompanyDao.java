package com.brian.jobs_db_service.dao;

import com.brian.jobs_db_service.model.entity.Company;

import java.math.BigInteger;
import java.util.Optional;

public interface CompanyDao {
    public Optional<Company> getCompanyById(Long id);
    public Optional<Company> getCompanyByName(String name);
    public Optional<Company> addCompany(String name);
}
