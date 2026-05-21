ALTER TABLE problem_contents
    DROP COLUMN IF EXISTS actual_graph,
    DROP COLUMN IF EXISTS expected_graph,
    DROP COLUMN IF EXISTS full_text,
    DROP COLUMN IF EXISTS author_id,
    ADD COLUMN IF NOT EXISTS description_md text NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS input_format_md text,
    ADD COLUMN IF NOT EXISTS output_format_md text,
    ADD COLUMN IF NOT EXISTS constraints_md text,
    ADD COLUMN IF NOT EXISTS notes_md text;

ALTER TABLE problem_contents
    ALTER COLUMN description_md DROP DEFAULT;
