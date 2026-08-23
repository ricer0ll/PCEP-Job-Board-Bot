package com.brian.jobs_db_service.dao;

import com.brian.jobs_db_service.model.entity.Company;
import com.brian.jobs_db_service.utils.Config;
import com.zaxxer.hikari.HikariConfig;
import com.zaxxer.hikari.HikariDataSource;
import org.jdbi.v3.core.Jdbi;
import org.jdbi.v3.core.mapper.reflect.BeanMapper;
import org.jdbi.v3.core.statement.Query;
import org.springframework.stereotype.Repository;

import java.math.BigInteger;
import java.util.Optional;

@Repository
public class CompanyDaoImpl implements CompanyDao {
    private static final String JDBI_URI_PREFIX = "jdbc:postgresql://";
    Jdbi jdbi;

    public CompanyDaoImpl() throws Exception {
        HikariConfig hikariConfig = new HikariConfig();
        hikariConfig.setJdbcUrl(formatUri(Config.getDbUri()));
        hikariConfig.setUsername(Config.getDbUser());
        hikariConfig.setPassword(Config.getDbPass());

        HikariDataSource ds = new HikariDataSource(hikariConfig);
        jdbi = Jdbi.create(ds);
    }

    private String formatUri(String uri) {
        return JDBI_URI_PREFIX + uri;
    }

    @Override
    public Optional<Company> getCompanyById(Long id) {
        return jdbi.withHandle(handle -> {
            Query query = handle
                    .createQuery("select id, name from public.companies where id = :company_id;")
                    .bind("company_id", id)
                    .registerRowMapper(BeanMapper.factory(Company.class));
            return query.mapTo(Company.class).findOne();
        });
    }

    @Override
    public Optional<Company> getCompanyByName(String name) {
        return jdbi.withHandle(handle -> {
            Query query = handle
                    .createQuery("select id, name from public.companies where name = :name;")
                    .bind("name", name)
                    .registerRowMapper(BeanMapper.factory(Company.class));
            return query.mapTo(Company.class).findOne();
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
