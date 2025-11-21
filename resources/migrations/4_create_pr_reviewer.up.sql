CREATE TABLE IF NOT EXISTS public.pr_reviewer
(
    id                  UUID                PRIMARY KEY         DEFAULT uuid_generate_v4(),
    pull_request_id     UUID                NOT NULL,
    reviewer_id         UUID                NOT NULL,
    assigned_at         TIMESTAMP           NOT NULL            DEFAULT NOW(),
    
    CONSTRAINT fk_pr_reviewer_pr 
        FOREIGN KEY (pull_request_id) 
        REFERENCES public.pull_request(pull_request_id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_pr_reviewer_user 
        FOREIGN KEY (reviewer_id) 
        REFERENCES public."user"(user_id) 
        ON DELETE CASCADE,

    CONSTRAINT uq_pr_reviewer UNIQUE (pull_request_id, reviewer_id)
);

CREATE INDEX IF NOT EXISTS idx_pr_reviewer_pr ON public.pr_reviewer(pull_request_id);
CREATE INDEX IF NOT EXISTS idx_pr_reviewer_user ON public.pr_reviewer(reviewer_id);
CREATE INDEX IF NOT EXISTS idx_pr_reviewer_assigned ON public.pr_reviewer(assigned_at DESC);

CREATE OR REPLACE FUNCTION check_max_reviewers()
RETURNS TRIGGER AS $$
BEGIN
    IF (SELECT COUNT(*) FROM public.pr_reviewer WHERE pull_request_id = NEW.pull_request_id) >= 2 THEN
        RAISE EXCEPTION 'Cannot assign more than 2 reviewers to a pull request';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_max_reviewers
    BEFORE INSERT ON public.pr_reviewer
    FOR EACH ROW
    EXECUTE FUNCTION check_max_reviewers();
