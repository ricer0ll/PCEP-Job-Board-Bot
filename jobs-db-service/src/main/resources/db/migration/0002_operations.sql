-- =========================================================
-- COMPANIES
-- =========================================================

-- Add a single company by name.
-- Returns the inserted row.
CREATE OR REPLACE FUNCTION public.add_company(
    p_name TEXT
)
RETURNS public.companies
LANGUAGE plpgsql
AS $$
DECLARE
    v_row public.companies;
BEGIN
    INSERT INTO public.companies (name)
    VALUES (p_name)
    ON CONFLICT (name) DO NOTHING
    RETURNING * INTO v_row;

    RETURN v_row;
END;
$$;

-- Get a single company by name.
-- Returns the matching row, or NULL if not found.
CREATE OR REPLACE FUNCTION public.get_company_by_name(
    p_name TEXT
)
RETURNS public.companies
LANGUAGE plpgsql
AS $$
DECLARE
    v_row public.companies;
BEGIN
    SELECT *
    INTO v_row
    FROM public.companies
    WHERE name = p_name;

    RETURN v_row; -- NULL (all fields) if no match found
END;
$$;


-- =========================================================
-- JOBS
-- =========================================================

-- Add a single job by id (and its company_id).
-- Returns the inserted row.
CREATE OR REPLACE FUNCTION public.add_job(
    p_id TEXT,
    p_company_id BIGINT
)
RETURNS public.jobs
LANGUAGE plpgsql
AS $$
DECLARE
    v_row public.jobs;
BEGIN
    INSERT INTO public.jobs (id, company_id)
    VALUES (p_id, p_company_id)
    RETURNING * INTO v_row;

    RETURN v_row;
END;
$$;

-- Get a single job by id and company_id.
-- Returns the matching row, or NULL if not found.
CREATE OR REPLACE FUNCTION public.get_job(
    p_id TEXT,
    p_company_id BIGINT
)
RETURNS public.jobs
LANGUAGE plpgsql
AS $$
DECLARE
    v_row public.jobs;
BEGIN
    SELECT *
    INTO v_row
    FROM public.jobs
    WHERE id = p_id
      AND company_id = p_company_id;

    RETURN v_row;
END;
$$;