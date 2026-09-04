package com.brian.jobs_db_service.utils;

import com.google.common.hash.Hashing;
import org.springframework.stereotype.Component;

import java.nio.charset.StandardCharsets;

@Component
public class HasherUtil {

    public static String getJobId(String jobTitle, String companyName) {
        return Hashing
                .sha256()
                .hashString(jobTitle + companyName, StandardCharsets.UTF_8)
                .toString();
    }
}
