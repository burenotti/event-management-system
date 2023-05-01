BEGIN;

CREATE TABLE events
(
    event_id            int8                     NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name                varchar(256)             NOT NULL,
    description         TEXT                     NOT NULL,
    organization_id     int8                     NOT NULL REFERENCES organizations ON DELETE CASCADE,
    creator_id          int8                     NOT NULL REFERENCES users ON DELETE SET NULL,
    registration_needed bool                     NOT NULL DEFAULT false,
    registration_begin  TIMESTAMP WITH TIME ZONE NULL     DEFAULT NULL,
    registration_end    TIMESTAMP WITH TIME ZONE NULL     DEFAULT NULL,
    begins_at           TIMESTAMP WITH TIME ZONE NOT NULL,
    ends_at             TIMESTAMP WITH TIME ZONE NULL     DEFAULT NULL,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    published_at        TIMESTAMP WITH TIME ZONE NULL     DEFAULT NULL
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_events_name ON events USING gin (name gin_trgm_ops) WITH (fastupdate = false);
CREATE INDEX idx_events_name_description ON events USING gin (description gin_trgm_ops) WITH (fastupdate = true);

COMMIT;