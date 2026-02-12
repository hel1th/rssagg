-- name: CreatePost :exec
INSERT INTO posts (
                  id,
                  created_at,
                  updated_at,
                  title,
                  description,
                  url,
                  feed_id,
                  published_at
                )
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (url) DO NOTHING;
-- RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*
FROM posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2 OFFSET $3;

-- -- name: GetNextFeedsToFetch :many
-- SELECT * FROM feeds
-- ORDER BY last_fetched_at NULLS FIRST
-- LIMIT $1;

-- -- name: MarkFeedAsFetched :one
-- UPDATE feeds
-- SET last_fetched_at = NOW(), updated_at = NOW()
-- WHERE id = $1
-- RETURNING *;
