package com.brian.jobs_db_service.utils;

import io.github.cdimascio.dotenv.Dotenv;

public final class Config {
    public static String getDbUser() {
        String user = System.getenv("POSTGRES_USER");
        if (user == null) {
            Dotenv dotenv = Dotenv.configure().load();
            return dotenv.get("POSTGRES_USER");
        }

        return user;
    }

    public static String getDbPass() {
        String pass = System.getenv("POSTGRES_PASSWORD");
        if (pass == null) {
            Dotenv dotenv = Dotenv.configure().load();
            return dotenv.get("POSTGRES_PASSWORD");
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
