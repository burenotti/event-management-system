BEGIN;


CREATE TABLE users
(
    user_id     int8        NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    first_name  varchar(32) NOT NULL,
    last_name   varchar(32) NOT NULL,
    middle_name varchar(32) NOT NULL     DEFAULT '',
    email       varchar(64) NOT NULL,
    is_active   bool        NOT NULL     DEFAULT FALSE,
    created_at  timestamp WITH TIME ZONE DEFAULT now()
);


COMMIT;