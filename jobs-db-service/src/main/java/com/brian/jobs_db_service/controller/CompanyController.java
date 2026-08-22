package com.brian.jobs_db_service.controller;

import com.brian.jobs_db_service.dao.CompanyDao;
import com.brian.jobs_db_service.dao.CompanyDaoImpl;
import com.brian.jobs_db_service.model.dto.company.AddCompanyRequest;
import com.brian.jobs_db_service.model.dto.company.AddCompanyResponse;
import com.brian.jobs_db_service.model.dto.company.GetCompanyResponse;
import com.brian.jobs_db_service.model.entity.Company;
import com.brian.jobs_db_service.utils.Config;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/v1/company")
public class CompanyController {
    private CompanyDao companyDao;

    public CompanyController() {
        try {
            companyDao = new CompanyDaoImpl(Config.getDbUri(), Config.getDbUser(), Config.getDbPass());
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    @PutMapping
    public AddCompanyResponse addCompany(@RequestBody AddCompanyRequest request) {
        Company company = companyDao.addCompany(request.getName());

        return new AddCompanyResponse(company.getId(), company.getName());
    }

    @GetMapping("/{company_id}")
    public GetCompanyResponse getCompanyById(@PathVariable Long company_id) {
        Company company = companyDao.getCompanyById(company_id);

        return new GetCompanyResponse(company.getId(), company.getName());
    }
}
