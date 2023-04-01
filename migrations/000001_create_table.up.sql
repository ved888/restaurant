BEGIN ;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS address
(
    id uuid DEFAULT uuid_generate_v4() primary key ,
    line1 text,
    line2 text ,
    pin_code text ,
    city text ,
    state text ,
    country text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);

CREATE TABLE IF NOT EXISTS users
(
    id uuid DEFAULT uuid_generate_v4 () primary key,
    first_name text ,
    middle_name text,
    last_name text,
    phone text ,
    email_id text ,
    city text ,
    address_id uuid references address(id) not null ,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone ,
    deleted_at timestamp without time zone

);

CREATE TABLE IF NOT EXISTS orders
(
    id uuid DEFAULT uuid_generate_v4 () primary key ,
    users_id uuid references users(id) not null ,
    status bigint,
    item_discount bigint,
    tax bigint,
    shipping text ,
    total bigint,
    deleted_at timestamp without time zone DEFAULT now(),
    created_at timestamp without time zone ,
    updated_at timestamp without time zone
);

CREATE TABLE IF NOT EXISTS billing
(
    id uuid DEFAULT uuid_generate_v4() primary key,
    users_id uuid references users(id) not null ,
    orders_id uuid references orders(id) not null ,
    type bigint,
    mode text,
    status bigint,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone ,
    deleted_at timestamp without time zone
);

CREATE TABLE IF NOT EXISTS res_table
(
    id uuid DEFAULT uuid_generate_v4 () primary key,
    code bigint,
    status bigint,
    capacity bigint,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone ,
    deleted_at timestamp without time zone
);

CREATE TABLE IF NOT EXISTS booking
(
    id uuid DEFAULT uuid_generate_v4 () primary key,
    booking_date date,
    rest_table_id uuid references res_table(id) not null ,
    users_id uuid references users(id) not null ,
    status bigint,
    pre_advance_booking boolean,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone ,
    deleted_at timestamp without time zone

);

CREATE TABLE IF NOT EXISTS food
(
    id uuid DEFAULT uuid_generate_v4 () primary key,
    name text ,
    price bigint,
    type text ,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);

CREATE TABLE IF NOT EXISTS order_item
(
    id uuid DEFAULT uuid_generate_v4 () primary key,
    orders_id uuid references orders(id) not null ,
    food_id uuid references food(id) not null ,
    price bigint,
    quantity bigint,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone ,
    deleted_at timestamp without time zone
);

COMMIT ;












