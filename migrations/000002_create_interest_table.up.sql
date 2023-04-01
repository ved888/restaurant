BEGIN ;

CREATE TABLE IF NOT EXISTS interest
(
    id uuid DEFAULT uuid_generate_v4 () primary key,
    name text ,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone ,
    deleted_at timestamp without time zone
);

CREATE TABLE IF NOT EXISTS relation_table
(
    id uuid DEFAULT uuid_generate_v4 () primary key ,
    users_id uuid references users(id) not null ,
    interest_id uuid references interest(id) not null,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone ,
    deleted_at timestamp without time zone
);

COMMIT ;