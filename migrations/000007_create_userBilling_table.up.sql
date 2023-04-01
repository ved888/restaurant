BEGIN;

ALTER TABLE billing
DROP COLUMN users_id;

CREATE TABLE IF NOT EXISTS user_billing(
    id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    users_id uuid references users(id) not null,
    billing_id uuid references billing(id) not null,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
    );

COMMIT;