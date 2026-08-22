package com.brian.jobs_db_service.dao;

import com.brian.jobs_db_service.model.entity.Company;
import com.zaxxer.hikari.HikariConfig;
import com.zaxxer.hikari.HikariDataSource;
import org.jdbi.v3.core.Jdbi;
import org.jdbi.v3.core.mapper.reflect.BeanMapper;
import org.jdbi.v3.core.statement.Query;

import java.math.BigInteger;
import java.util.Optional;

public class CompanyDaoImpl implements CompanyDao {
    private static final String JDBI_URI_PREFIX = "jdbc:postgresql://";
    Jdbi jdbi;

    public CompanyDaoImpl(String uri, String user, String pass) throws Exception {
        HikariConfig config = new HikariConfig();
        config.setJdbcUrl(formatUri(uri));
        config.setUsername(user);
        config.setPassword(pass);

        HikariDataSource ds = new HikariDataSource(config);
        jdbi = Jdbi.create(ds);
    }

    private String formatUri(String uri) {
        return JDBI_URI_PREFIX + uri;
    }

    @Override
    public Company getCompanyById(Long id) {
        return jdbi.withHandle(handle -> {
            Query query = handle
                    .createQuery("select id, name from public.companies where id = :company_id;")
                    .bind("company_id", id)
                    .registerRowMapper(BeanMapper.factory(Company.class));
            return query.mapTo(Company.class).one();
        });
    }

    @Override
    public Company addCompany(String name) {
        return jdbi.withHandle(handle -> {
            Query query = handle
                    .createQuery("select * from public.add_company(:company_name);")
                    .bind("company_name", name)
                    .registerRowMapper(BeanMapper.factory(Company.class));
            return query.mapTo(Company.class).one();
        });
    }
}
