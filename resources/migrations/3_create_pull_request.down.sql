DROP INDEX IF EXISTS idx_pr_created;
DROP INDEX IF EXISTS idx_pr_status;
DROP INDEX IF EXISTS idx_pr_author;

DROP TABLE IF EXISTS public.pull_request CASCADE;

DROP TYPE IF EXISTS pr_status;
