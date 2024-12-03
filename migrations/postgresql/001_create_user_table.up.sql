CREATE TABLE public.user (
    id UUID PRIMARY KEY NOT NULL,
    refresh_token TEXT,
    email VARCHAR(255)
)