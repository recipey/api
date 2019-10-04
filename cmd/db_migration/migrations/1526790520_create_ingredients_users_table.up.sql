-- After some reading looks like for larger applications a trick done for scalability is to not
-- use Foreign Keys (FK) at all. To enforce FK constraints the RDMS will have to check every row
-- to be sure; this impacts INSERT, UPDATE queries.
CREATE TABLE IF NOT EXISTS ingredients_users (
  ingredient_id INT,
  user_id INT,
  quantity INT,
  unit TEXT,
  --  FOREIGN KEY (user_id) REFERENCES users (id),
  PRIMARY KEY (ingredient_id, user_id)
);
