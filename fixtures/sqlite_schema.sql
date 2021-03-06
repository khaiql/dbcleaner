CREATE TABLE users (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name TEXT
);

CREATE TABLE addresses (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER null,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
