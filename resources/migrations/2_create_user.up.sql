CREATE TABLE IF NOT EXISTS public."user"
(
    user_id         UUID                PRIMARY KEY         DEFAULT uuid_generate_v4(),
    username        VARCHAR(255)        NOT NULL,
    team_id         UUID                NOT NULL,
    is_active       BOOLEAN             NOT NULL            DEFAULT true,
    created_at      TIMESTAMP           NOT NULL            DEFAULT NOW(),
    updated_at      TIMESTAMP           NOT NULL            DEFAULT NOW(),
    
    CONSTRAINT fk_user_team 
        FOREIGN KEY (team_id) 
        REFERENCES public.team(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_team_id ON public."user"(team_id);
CREATE INDEX IF NOT EXISTS idx_user_active ON public."user"(is_active);
CREATE INDEX IF NOT EXISTS idx_user_username ON public."user"(username);
