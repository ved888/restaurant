BEGIN;

ALTER TABLE booking
DROP COLUMN rest_table_id;

CREATE TABLE IF NOT EXISTS booking_table(
    id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    booking_id uuid references booking(id) not null,
    rest_table_id uuid references res_table(id) not null,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
    );

COMMIT;