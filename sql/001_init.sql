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

BEGIN;
CREATE TABLE IF NOT EXISTS public.user
(
    account_id VARCHAR(32) primary key,
    email_address TEXT not null,
    login_password_hash TEXT not null,  /* pgcryptoのcrypt関数を使用 */
    name TEXT,
    belongings TEXT,
    is_admin boolean not null default 0,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
WITH (
    OIDS = FALSE
);

COMMENT ON TABLE public.user
    IS 'Store management-api users.';
END;


