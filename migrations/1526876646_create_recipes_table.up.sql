CREATE TABLE IF NOT EXISTS recipes (
  id SERIAL PRIMARY KEY NOT NULL,
  name TEXT NOT NULL,
  author TEXT,
  image_url TEXT,
  source_url TEXT
);
