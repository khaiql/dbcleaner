CREATE TABLE users (
  id bigint not null AUTO_INCREMENT,
  primary key (id)
) ENGINE=INNODB;

CREATE TABLE addresses (
  id bigint not null  AUTO_INCREMENT,
  user_id bigint null,
  primary key (id),
  FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=INNODB;
