BEGIN;

CREATE TABLE login_code
(
    user_id    int8                     NOT NULL REFERENCES users,
    code       varchar(6)               NOT NULL,
    is_used    bool                     NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_login_code ON login_code (user_id, code);

COMMIT;