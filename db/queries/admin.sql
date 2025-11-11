-- Admin dashboard queries

-- Get total user count
-- name: GetTotalUserCount :one
SELECT COUNT(*) FROM users;

-- Get user count for today
-- name: GetUserCountToday :one
SELECT COUNT(*) FROM users WHERE DATE(created_at) = CURRENT_DATE;

-- Get user count for this week
-- name: GetUserCountThisWeek :one
SELECT COUNT(*) FROM users WHERE created_at >= DATE_TRUNC('week', CURRENT_DATE);

-- Get recent user activity (last 10 registrations)
-- name: GetRecentUsers :many
SELECT id, email, name, created_at FROM users ORDER BY created_at DESC LIMIT 10;

-- Get admin users
-- name: GetAdminUsers :many
SELECT id, email, name, created_at FROM users WHERE is_admin = TRUE ORDER BY created_at DESC;