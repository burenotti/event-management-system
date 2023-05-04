BEGIN;

ALTER TABLE organization_members
    DROP COLUMN can_view_events;
ALTER TABLE organization_members
    ADD COLUMN is_owner bool NOT NULL DEFAULT FALSE;

COMMIT;