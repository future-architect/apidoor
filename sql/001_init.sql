BEGIN;
CREATE EXTENSION IF NOT EXISTS pgcrypto;


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

CREATE TABLE IF NOT EXISTS public.product
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL unique,
    /*
    TODO: ownerの設定 https://github.com/future-architect/apidoor/issues/79
    owner INT REFERENCES apiuser(id),
     */
    display_name TEXT NOT NULL,
    source TEXT NOT NULL,
    description TEXT NOT NULL,
    thumbnail TEXT NOT NULL,
    base_path TEXT NOT NULL,
    swagger_url TEXT NOT NULL,
    is_available SMALLINT DEFAULT 0 NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
    WITH (
        OIDS = FALSE
    );


COMMENT ON TABLE public.product
    IS 'Store information of products(API).';

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

CREATE TABLE IF NOT EXISTS public.contract
(
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES apiuser(id),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.contract_product_content
(
    id SERIAL PRIMARY KEY,
    contract_id INT REFERENCES contract(id),
    product_id INT REFERENCES product(id),
    description TEXT,
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

CREATE TABLE IF NOT EXISTS public.apikey_contract_product_authorized
(
    id SERIAL PRIMARY KEY,
    apikey_id INT REFERENCES apikey(id),
    contract_product_id INT REFERENCES contract_product_content(id),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);


END;
