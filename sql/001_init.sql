BEGIN;
CREATE TABLE IF NOT EXISTS public.apiinfo
(
    id serial NOT NULL,
    name text NOT NULL,
    source text NOT NULL,
    description text NOT NULL,
    thumbnail text NOT NULL,
    swagger_url text NOT NULL
)
WITH (
    OIDS = FALSE
);

COMMENT ON TABLE public.apiinfo
    IS 'Store information of products(API).';
END;
