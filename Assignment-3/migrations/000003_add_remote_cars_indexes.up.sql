CREATE INDEX IF NOT EXISTS remote_cars_title_idx ON remote_cars USING GIN (to_tsvector('simple', name));
