CREATE TABLE IF NOT EXISTS checkpoints (
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL,
  slug varchar(255) NOT NULL UNIQUE,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW()
);
