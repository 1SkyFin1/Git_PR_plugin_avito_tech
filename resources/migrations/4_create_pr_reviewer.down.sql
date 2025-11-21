DROP TRIGGER IF EXISTS trg_check_max_reviewers ON public.pr_reviewer;
DROP FUNCTION IF EXISTS check_max_reviewers();

DROP INDEX IF EXISTS idx_pr_reviewer_assigned;
DROP INDEX IF EXISTS idx_pr_reviewer_user;
DROP INDEX IF EXISTS idx_pr_reviewer_pr;

DROP TABLE IF EXISTS public.pr_reviewer CASCADE;
