CREATE TABLE users (
  id SERIAL,
  name varchar,
  primary key (id)
);

CREATE TABLE addresses (
  id SERIAL,
  user_id bigint null,
  primary key (id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);
