CREATE TABLE IF NOT EXISTS problems (
    id uuid PRIMARY KEY,
    name varchar(255) NOT NULL,
    slug varchar(255) NOT NULL,
    difficulty varchar(32) NOT NULL,
    tag varchar(64),
    status varchar(20) NOT NULL DEFAULT 'draft',
    author_id uuid NOT NULL,
    verified_at timestamptz NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    CONSTRAINT problems_status_check CHECK (status IN ('draft', 'approved', 'rejected')),
    CONSTRAINT uni_problems_slug UNIQUE (slug)  -- явное имя constraint
);

CREATE INDEX IF NOT EXISTS idx_problems_difficulty ON problems (difficulty);
CREATE INDEX IF NOT EXISTS idx_problems_tag ON problems (tag);
CREATE INDEX IF NOT EXISTS idx_problems_status ON problems (status);
CREATE INDEX IF NOT EXISTS idx_problems_author_id ON problems (author_id);
CREATE INDEX IF NOT EXISTS idx_problems_verified_at ON problems (verified_at);
CREATE INDEX IF NOT EXISTS idx_problems_created_at ON problems (created_at DESC);

CREATE TABLE IF NOT EXISTS problem_contents (
    id uuid PRIMARY KEY,
    problem_id uuid NOT NULL,
    description_md text NOT NULL,
    input_format_md text,
    output_format_md text,
    constraints_md text,
    notes_md text,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    CONSTRAINT uni_problem_contents_problem_id UNIQUE (problem_id),  -- добавляем UNIQUE constraint
    CONSTRAINT fk_problem_contents_problem_id
        FOREIGN KEY (problem_id)
        REFERENCES problems (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_problem_contents_problem_id ON problem_contents (problem_id);

