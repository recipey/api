CREATE TABLE IF NOT EXISTS recipe_steps (
  recipe_id INTEGER,
  step_no INTEGER,
  ingredient_id INTEGER,
  ingredient_unit TEXT,
  ingredient_quantity REAL,
  description TEXT NOT NULL,
  FOREIGN KEY (recipe_id) REFERENCES recipes (id),
  PRIMARY KEY (recipe_id, step_no)
)
