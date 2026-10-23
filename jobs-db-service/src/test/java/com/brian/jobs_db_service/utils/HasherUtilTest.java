package com.brian.jobs_db_service.utils;

import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class HasherUtilTest {

    @Test
    void getJobId_returnsHexSha256OfTitleConcatenatedWithCompany() {
        assertThat(HasherUtil.getJobId("Engineer", "Acme"))
                .isEqualTo(HasherUtil.getJobId("EngineerAcme", ""))
                .hasSize(64)
                .matches("[0-9a-f]{64}");
    }

    @Test
    void getJobId_isDeterministic() {
        assertThat(HasherUtil.getJobId("Engineer", "Acme"))
                .isEqualTo(HasherUtil.getJobId("Engineer", "Acme"));
    }

    @Test
    void getJobId_matchesKnownHash() {
        // echo -n "abc" | sha256sum
        assertThat(HasherUtil.getJobId("a", "bc"))
                .isEqualTo("ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad");
    }

    @Test
    void getJobId_differsForDifferentInputs() {
        assertThat(HasherUtil.getJobId("Engineer", "Acme"))
                .isNotEqualTo(HasherUtil.getJobId("Engineer", "Globex"))
                .isNotEqualTo(HasherUtil.getJobId("Manager", "Acme"));
    }
}
