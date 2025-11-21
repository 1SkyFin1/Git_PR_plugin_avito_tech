DROP INDEX IF EXISTS idx_user_username;
DROP INDEX IF EXISTS idx_user_active;
DROP INDEX IF EXISTS idx_user_team_id;

DROP TABLE IF EXISTS public."user" CASCADE;
