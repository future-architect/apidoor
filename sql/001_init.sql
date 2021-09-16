BEGIN;
CREATE TABLE IF NOT EXISTS public.apiinfo
(
    id integer NOT NULL,
    name text NOT NULL,
    source text NOT NULL,
    description text NOT NULL,
    thumbnail text NOT NULL
)
WITH (
    OIDS = FALSE
);

COMMENT ON TABLE public.apiinfo
    IS 'Store information of products(API).';

CREATE TABLE IF NOT EXISTS public.log_list
(
    run_date timestamp with time zone NOT NULL,
    api_key text NOT NULL,
    api_path text NOT NULL,
    custom_log jsonb
);

COMMENT ON TABLE public.log_list
    IS 'Table to save log of gateway.';
END;