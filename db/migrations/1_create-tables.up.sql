CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS customers (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text,
    balance float default 0,
    created_at TIMESTAMP WITHOUT TIME ZONE default CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE IF NOT EXISTS payments (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id uuid NOT NULL,
    payment_date text,
    amount float default 0,
    verified boolean default false,
    created_at TIMESTAMP WITHOUT TIME ZONE default CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);

CREATE TABLE IF NOT EXISTS transactions (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id uuid NOT NULL,
    shipment_date text NOT NULL,
    shipment_cost float default 0,
    total float default 0,
    created_at TIMESTAMP WITHOUT TIME ZONE default CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);

CREATE TABLE IF NOT EXISTS items (
    item text PRIMARY KEY NOT NULL,
);

CREATE TABLE IF NOT EXISTS transaction_items (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id uuid NOT NULL,
    item text NOT NULL,
    current_price float default 0,
    kilos float default 0,
    created_at TIMESTAMP WITHOUT TIME ZONE default CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);

CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    username text NOT NULL,
    encrypted_password text NOT NULL,
    user_role text,
    customer_id text,
    created_at TIMESTAMP WITHOUT TIME ZONE default CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

