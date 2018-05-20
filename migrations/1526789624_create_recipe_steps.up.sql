CREATE TABLE IF NOT EXISTS recipe_steps (
  id serial int NOT NULL,
  step_no int NOT NULL,
  PRIMARY KEY (id, step_no),
  ingredient_unit TEXT,
  ingredient_quantity TEXT,
  description TEXT NOT NULL,
  CHECK (description <> '')
);
