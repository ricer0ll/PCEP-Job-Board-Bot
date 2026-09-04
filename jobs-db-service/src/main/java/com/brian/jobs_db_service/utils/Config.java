package com.brian.jobs_db_service.utils;

import com.brian.jobs_db_service.model.dto.job.GetJobIdRequest;
import com.brian.jobs_db_service.model.dto.job.GetJobIdResponse;
import com.google.common.hash.Hashing;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Component;
import org.springframework.web.bind.annotation.RequestBody;

import java.nio.charset.StandardCharsets;

@Component
public class Config {
    public static String getDbUser() {
        String user = System.getenv("POSTGRES_USER");
        if (user == null) {
            throw new RuntimeException("POSTGRES_USER environment variable not found");
        }
        return user;
    }

    public static String getDbPass() throws Exception {
        String pass = System.getenv("POSTGRES_PASSWORD");
        if (pass == null) {
            throw new RuntimeException("POSTGRES_PASSWORD environment variable not found");
        }
        return pass;
    }

    public static String getDbUri() {
        String env = System.getenv("ENV");

        if (env == null) {
            return "localhost:5432/postgres";
        }

        if (env.equalsIgnoreCase("dev")) {
            return "localhost:5432/postgres";
        } else {
            return "db:5432/postgres";
        }
    }
}
