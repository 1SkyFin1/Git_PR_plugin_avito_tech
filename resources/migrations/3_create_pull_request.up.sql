CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE IF NOT EXISTS public.pull_request
(
    pull_request_id     UUID                PRIMARY KEY         DEFAULT uuid_generate_v4(),
    pull_request_name   VARCHAR(500)        NOT NULL,
    author_id           UUID                NOT NULL,
    status              pr_status           NOT NULL            DEFAULT 'OPEN',
    created_at          TIMESTAMP           NOT NULL            DEFAULT NOW(),
    merged_at           TIMESTAMP           NULL,
    
    CONSTRAINT fk_pr_author 
        FOREIGN KEY (author_id) 
        REFERENCES public."user"(user_id) 
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_pr_author ON public.pull_request(author_id);
CREATE INDEX IF NOT EXISTS idx_pr_status ON public.pull_request(status);
CREATE INDEX IF NOT EXISTS idx_pr_created ON public.pull_request(created_at DESC);
