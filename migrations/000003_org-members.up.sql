BEGIN;

CREATE TABLE organization_members
(
    organization_id    int8 NOT NULL REFERENCES organizations ON DELETE CASCADE,
    user_id            int8 NOT NULL REFERENCES users ON DELETE CASCADE,
    can_edit_events    bool NOT NULL DEFAULT FALSE,
    can_manage_members bool NOT NULL DEFAULT FALSE,
    can_view_events    bool NOT NULL DEFAULT TRUE,
    PRIMARY KEY (organization_id, user_id)
);

COMMIT;