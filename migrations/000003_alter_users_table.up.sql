BEGIN;

ALTER TABLE users
DROP COLUMN address_id ;

    CREATE TABLE IF NOT EXISTS user_address(
        id uuid DEFAULT uuid_generate_v4 () primary key ,
        user_id uuid references users(id) not null ,
        address_id uuid references address(id) not null,
        created_at timestamp without time zone DEFAULT now(),
        updated_at timestamp without time zone ,
        deleted_at timestamp without time zone
    );

COMMIT ;