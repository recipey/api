CREATE TABLE IF NOT EXISTS recipe_steps (
  recipe_id INT,
  step_no INT,
  ingredient_id INT,
  ingredient_unit TEXT,
  ingredient_quantity FLOAT,
  description TEXT NOT NULL,
  FOREIGN KEY (recipe_id) REFERENCES recipes (id),
  PRIMARY KEY (recipe_id, step_no)
)
