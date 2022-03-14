BEGIN;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS public.apiinfo
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    source TEXT NOT NULL,
    description TEXT NOT NULL,
    thumbnail TEXT NOT NULL,
    swagger_url TEXT NOT NULL
)
WITH (
    OIDS = FALSE
);

COMMENT ON TABLE public.apiinfo
    IS 'Store information of products(API).';

CREATE TABLE IF NOT EXISTS public.apiuser
(
    id SERIAL PRIMARY KEY ,
    account_id VARCHAR(32) NOT NULL UNIQUE ,
    email_address TEXT NOT NULL,
    login_password_hash TEXT NOT NULL,  /* pgcryptoのcrypt関数を使用 */
    name TEXT,
    permission_flag VARCHAR(2) NOT NULL DEFAULT '00',
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
WITH (
    OIDS = FALSE
);

COMMENT ON TABLE public.apiuser
    IS 'Store management-api users.';

CREATE TABLE IF NOT EXISTS public.log_list
(
    id SERIAL PRIMARY KEY ,
    run_date TIMESTAMP WITH TIME ZONE NOT NULL,
    api_key TEXT NOT NULL,
    api_path TEXT NOT NULL,
    custom_log jsonb
);

COMMENT ON TABLE public.log_list
    IS 'Table to save log of gateway.';

CREATE TABLE IF NOT EXISTS public.product
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL unique,
    display_name TEXT,
    source TEXT NOT NULL,
    description TEXT NOT NULL,
    thumbnail TEXT NOT NULL,
    is_available SMALLINT DEFAULT 0 NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.product_api_content
(
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES product(id),
    api_id INT REFERENCES apiinfo(id),
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.contract
(
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES apiuser(id),
    product_id INT REFERENCES product(id),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.apikey
(
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES apiuser(id),
    access_key TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

END;
