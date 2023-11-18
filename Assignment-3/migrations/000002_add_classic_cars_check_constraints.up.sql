ALTER TABLE classic_cars ADD CONSTRAINT classic_cars_cost_check CHECK (cost > 0);
ALTER TABLE classic_cars ADD CONSTRAINT classic_cars_year_check CHECK (year BETWEEN 1 AND date_part('year', now()));
