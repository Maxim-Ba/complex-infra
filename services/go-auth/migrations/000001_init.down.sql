DROP TRIGGER IF EXISTS update_users_timestamp ON users;
DROP FUNCTION IF EXISTS update_timestamp;
DROP INDEX IF EXISTS idx_users_login;
DROP TABLE IF EXISTS users;
