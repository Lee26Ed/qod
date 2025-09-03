CREATE TABLE IF NOT EXISTS quotes (
    id bigserial PRIMARY KEY,
    content TEXT NOT NULL,
    author TEXT NOT NULL,
    version integer NOT NULL DEFAULT 1,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);