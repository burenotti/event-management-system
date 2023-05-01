BEGIN;

CREATE TABLE organizations
(
    organization_id int8         NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name            varchar(256) NOT NULL,
    address         varchar(256) NULL DEFAULT NULL,
    contact_email   varchar(64)  NULL DEFAULT NULL,
    contact_phone   varchar(32)  NULL DEFAULT NULL
);


COMMIT;