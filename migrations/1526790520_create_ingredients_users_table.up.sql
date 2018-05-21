CREATE TABLE IF NOT EXISTS ingredients_users (
  ingredient_id INTEGER,
  user_id INTEGER,
  quantity INTEGER,
  unit TEXT,
  FOREIGN KEY (user_id) REFERENCES users (id),
  PRIMARY KEY (ingredient_id, user_id)
);
