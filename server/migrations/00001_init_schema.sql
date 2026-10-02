-- +goose Up
CREATE TABLE tasks (
    id             BIGSERIAL PRIMARY KEY,
    title          TEXT NOT NULL,
    status         TEXT NOT NULL DEFAULT 'open'
                   CHECK (status IN ('open', 'blocked', 'done', 'dropped')),
    blocker_reason TEXT,
    important      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    closed_at      TIMESTAMPTZ
);

CREATE TABLE events (
    id         BIGSERIAL PRIMARY KEY,
    task_id    BIGINT REFERENCES tasks (id),
    type       TEXT NOT NULL
               CHECK (type IN ('created', 'progress', 'blocked', 'unblocked', 'done', 'dropped', 'note')),
    text       TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX events_created_at_idx ON events (created_at);
CREATE INDEX events_task_id_idx ON events (task_id);

CREATE TABLE digests (
    digest_date  DATE PRIMARY KEY,
    summary      TEXT NOT NULL,
    plan         TEXT NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    sent_at      TIMESTAMPTZ
);

CREATE TABLE devices (
    id           BIGSERIAL PRIMARY KEY,
    push_token   TEXT NOT NULL UNIQUE,
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE devices;
DROP TABLE digests;
DROP TABLE events;
DROP TABLE tasks;
