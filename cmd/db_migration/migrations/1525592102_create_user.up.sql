-- After some reading looks like it's recommended to always use TEXT type over varchar(n) because
-- performance benefits are minimal if any.
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY NOT NULL,
  username TEXT UNIQUE,
  email TEXT UNIQUE,
  password TEXT
);
