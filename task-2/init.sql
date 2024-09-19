DROP TABLE IF EXISTS users;

CREATE TABLE users (
  users_id SERIAL  PRIMARY KEY,
  name VARCHAR NOT NULL
);

INSERT INTO users
  (name)
VALUES
  ('Justin Case'),
  ('Max Power');
