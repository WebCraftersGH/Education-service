ALTER TABLE problem_contents
    DROP COLUMN IF EXISTS description_md,
    DROP COLUMN IF EXISTS input_format_md,
    DROP COLUMN IF EXISTS output_format_md,
    DROP COLUMN IF EXISTS constraints_md,
    DROP COLUMN IF EXISTS notes_md,
    ADD COLUMN IF NOT EXISTS actual_graph jsonb NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS expected_graph jsonb NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS full_text text NOT NULL DEFAULT '';

ALTER TABLE problem_contents
    ALTER COLUMN actual_graph DROP DEFAULT,
    ALTER COLUMN expected_graph DROP DEFAULT,
    ALTER COLUMN full_text DROP DEFAULT;

ALTER TABLE problem_contents
    ADD COLUMN IF NOT EXISTS author_id uuid;

UPDATE problem_contents AS pc
SET author_id = p.author_id
FROM problems AS p
WHERE pc.problem_id = p.id
  AND pc.author_id IS NULL;

ALTER TABLE problem_contents
    ALTER COLUMN author_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_problem_contents_author_id ON problem_contents (author_id);
