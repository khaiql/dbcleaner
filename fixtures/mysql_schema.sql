CREATE TABLE users (
  id bigint not null AUTO_INCREMENT,
  name varchar(255),
  primary key (id)
);

CREATE TABLE addresses (
  id bigint not null  AUTO_INCREMENT,
  user_id bigint null,
  primary key (id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);
