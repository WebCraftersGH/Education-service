CREATE TABLE IF NOT EXISTS checkpoints (
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL,
  slug varchar(255) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX checkpoints_user_id_slug_key ON checkpoints (user_id, slug);
