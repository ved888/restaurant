BEGIN;

ALTER TABLE booking
DROP COLUMN users_id;

CREATE TABLE IF NOT EXISTS user_booking(
    id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    users_id uuid references users(id) not null,
    booking_id uuid references booking(id) not null,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
    );

COMMIT;