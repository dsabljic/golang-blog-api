CREATE TABLE IF NOT EXISTS user_tokens (
  token bytea PRIMARY KEY,
  user_id bigint NOT NULL
)
