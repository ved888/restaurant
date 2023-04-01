BEGIN;

ALTER TABLE orders
DROP COLUMN users_id;

CREATE TABLE IF NOT EXISTS user_order(
    id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    users_id uuid references users(id) not null,
    orders_id uuid references orders(id) not null,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);

COMMIT;