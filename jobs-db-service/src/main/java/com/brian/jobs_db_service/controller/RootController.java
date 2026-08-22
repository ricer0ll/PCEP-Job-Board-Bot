package com.brian.jobs_db_service.controller;

import com.brian.jobs_db_service.model.dto.greeting.GreetingResponse;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/root")
public class RootController {
    private static final String greetingMessage = "Hello, World!";


    @GetMapping
    public GreetingResponse greeting() {
        return new GreetingResponse(greetingMessage);
    }
}
