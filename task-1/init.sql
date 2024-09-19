DROP TABLE IF EXISTS app_user;

CREATE TABLE app_user (
  app_user_id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL,
  email VARCHAR NOT NULL UNIQUE,
  age INTEGER NOT NULL
);

INSERT INTO app_user
  (name, email, age)
VALUES
  ('Justin Case', 'justin@case.com', 33);
