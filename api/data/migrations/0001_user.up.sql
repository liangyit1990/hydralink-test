CREATE TABLE public.users
(
    id              TEXT                     NOT NULL CHECK (id <> ''::TEXT),
    first_name      TEXT                     NOT NULL CHECK (first_name <> ''::TEXT),
    last_name       TEXT                     NOT NULL CHECK (last_name <> ''::TEXT),
    email           TEXT                     NOT NULL CHECK (email <> ''::TEXT), -- also username
    hashed_password TEXT                     NOT NULL CHECK (hashed_password <> ''::TEXT), -- hashed password
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
-- This ensures no duplicated id in table
CREATE UNIQUE INDEX user_index_id ON public.users (id);
-- This ensures no duplicated first & last name in table
CREATE UNIQUE INDEX user_index_first_last_name ON public.users (first_name,last_name);
-- This ensures no duplicated email in table
CREATE UNIQUE INDEX user_index_email ON public.users (email);
