CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.team
(
    id              UUID                PRIMARY KEY         DEFAULT uuid_generate_v4(),
    team_name       VARCHAR(255)        NOT NULL            UNIQUE,
    created_at      TIMESTAMP           NOT NULL            DEFAULT NOW(),
    updated_at      TIMESTAMP           NOT NULL            DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_team_name ON public.team(team_name);
