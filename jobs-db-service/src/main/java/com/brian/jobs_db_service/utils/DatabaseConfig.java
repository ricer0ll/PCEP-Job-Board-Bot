package com.brian.jobs_db_service.utils;

import com.zaxxer.hikari.HikariConfig;
import com.zaxxer.hikari.HikariDataSource;
import org.jdbi.v3.core.Jdbi;
import org.jdbi.v3.sqlobject.SqlObjectPlugin;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import javax.sql.DataSource;

@Configuration
public class DatabaseConfig {
    private static final String JDBI_URI_PREFIX = "jdbc:postgresql://";

    @Bean
    public DataSource dataSource() throws Exception {
        HikariConfig hikariConfig = new HikariConfig();
        hikariConfig.setJdbcUrl(JDBI_URI_PREFIX + Config.getDbUri());
        hikariConfig.setUsername(Config.getDbUser());
        hikariConfig.setPassword(Config.getDbPass());

        return new HikariDataSource(hikariConfig);
    }

    @Bean
    public Jdbi jdbi(DataSource dataSource) {
        Jdbi jdbi = Jdbi.create(dataSource);

        jdbi.installPlugin(new SqlObjectPlugin());

        return jdbi;
    }
}
