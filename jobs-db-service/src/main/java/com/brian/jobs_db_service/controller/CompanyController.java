package com.brian.jobs_db_service.controller;

import com.brian.jobs_db_service.dao.CompanyDao;
import com.brian.jobs_db_service.model.dto.company.AddCompanyRequest;
import com.brian.jobs_db_service.model.dto.company.AddCompanyResponse;
import com.brian.jobs_db_service.model.dto.company.GetCompanyResponse;
import com.brian.jobs_db_service.model.entity.Company;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

import java.util.Optional;

@RestController
@RequestMapping("/v1/company")
public class CompanyController {
    private final CompanyDao companyDao;

    public CompanyController(CompanyDao companyDao) {
        this.companyDao = companyDao;
    }

    @PostMapping
    public ResponseEntity<AddCompanyResponse> addCompany(@RequestBody AddCompanyRequest request) {
        Company company = companyDao.addCompany(request.getName())
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.CONFLICT, "Company already exists"));

        return ResponseEntity.ok(new AddCompanyResponse(company.getId(), company.getName()));
    }

    @GetMapping("/{company_id}")
    public ResponseEntity<GetCompanyResponse> getCompanyById(@PathVariable Long company_id) {
        Optional<Company> companyOpt = companyDao.getCompanyById(company_id);
        if (companyOpt.isEmpty()) {
            return ResponseEntity.notFound().build();
        }

        Company company = companyOpt.get();
        GetCompanyResponse response = new GetCompanyResponse(company.getId(), company.getName());
        return ResponseEntity.ok(response);
    }

    @GetMapping("/name/{company_name}")
    public ResponseEntity<GetCompanyResponse> getCompanyByName(@PathVariable String company_name) {
        Optional<Company> companyOpt = companyDao.getCompanyByName(company_name);
        if (companyOpt.isEmpty()) {
            return ResponseEntity.notFound().build();
        }

        Company company = companyOpt.get();
        GetCompanyResponse response = new GetCompanyResponse(company.getId(), company.getName());
        return ResponseEntity.ok(response);
    }
}
