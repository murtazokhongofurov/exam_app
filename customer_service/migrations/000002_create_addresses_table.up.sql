CREATE TABLE IF NOT EXISTS addresses(
    id UUID PRIMARY KEY NOT NULL,
    owner_id UUID NOT NULL REFERENCES customers(id),
    country VARCHAR(50) NOT NULL,
    street VARCHAR(50) NOT NULL
);