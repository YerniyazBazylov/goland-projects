CREATE INDEX IF NOT EXISTS classic_cars_title_idx ON classic_cars USING GIN (to_tsvector('simple', name));
