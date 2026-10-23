package com.brian.jobs_db_service.controller;

import com.brian.jobs_db_service.dao.CompanyDao;
import com.brian.jobs_db_service.model.entity.Company;
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

@WebMvcTest(CompanyController.class)
class CompanyControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockitoBean
    private CompanyDao companyDao;

    @Test
    void addCompany_returnsCreatedCompany() throws Exception {
        when(companyDao.addCompany("Acme")).thenReturn(Optional.of(new Company(1L, "Acme")));

        mockMvc.perform(post("/v1/company")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("{\"name\":\"Acme\"}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(1))
                .andExpect(jsonPath("$.name").value("Acme"));
    }

    @Test
    void addCompany_returnsConflictWhenCompanyExists() throws Exception {
        when(companyDao.addCompany("Acme")).thenReturn(Optional.empty());

        mockMvc.perform(post("/v1/company")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("{\"name\":\"Acme\"}"))
                .andExpect(status().isConflict());
    }

    @Test
    void getCompanyById_returnsCompany() throws Exception {
        when(companyDao.getCompanyById(1L)).thenReturn(Optional.of(new Company(1L, "Acme")));

        mockMvc.perform(get("/v1/company/1"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(1))
                .andExpect(jsonPath("$.name").value("Acme"));
    }

    @Test
    void getCompanyById_returnsNotFoundWhenMissing() throws Exception {
        when(companyDao.getCompanyById(1L)).thenReturn(Optional.empty());

        mockMvc.perform(get("/v1/company/1"))
                .andExpect(status().isNotFound());
    }

    @Test
    void getCompanyById_returnsBadRequestForNonNumericId() throws Exception {
        mockMvc.perform(get("/v1/company/abc"))
                .andExpect(status().isBadRequest());
    }

    @Test
    void getCompanyByName_returnsCompany() throws Exception {
        when(companyDao.getCompanyByName("Acme")).thenReturn(Optional.of(new Company(1L, "Acme")));

        mockMvc.perform(get("/v1/company/name/Acme"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(1))
                .andExpect(jsonPath("$.name").value("Acme"));
    }

    @Test
    void getCompanyByName_returnsNotFoundWhenMissing() throws Exception {
        when(companyDao.getCompanyByName("Acme")).thenReturn(Optional.empty());

        mockMvc.perform(get("/v1/company/name/Acme"))
                .andExpect(status().isNotFound());
    }
}
