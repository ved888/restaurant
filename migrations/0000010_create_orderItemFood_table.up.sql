BEGIN;

CREATE TABLE IF NOT EXISTS order_item_food(
    id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    food_id uuid references food(id) not null,
    order_item_id uuid references order_item(id) not null,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
    );

COMMIT;