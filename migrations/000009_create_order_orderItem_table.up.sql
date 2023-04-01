BEGIN;

ALTER TABLE order_item
DROP COLUMN orders_id;

CREATE TABLE IF NOT EXISTS order_orderItem(
    id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    order_id uuid references orders(id) not null,
    orderItem_id uuid references order_item(id) not null,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
    );

COMMIT;