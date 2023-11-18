CREATE TABLE IF NOT EXISTS classic_cars (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    year integer NOT NULL,
    cost integer NOT NULL,
    description text NOT NULL,
    version integer NOT NULL DEFAULT 1
    );
