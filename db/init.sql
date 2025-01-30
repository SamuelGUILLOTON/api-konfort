DROP TABLE IF EXISTS user;

CREATE TABLE user (
  id         INT AUTO_INCREMENT NOT NULL,
  blaze      VARCHAR(128) NOT NULL,
  password     VARCHAR(255) NOT NULL,
  email       VARCHAR(255) NOT NULL UNIQUE,
  role      VARCHAR(255) NOT NULL,
  created_at DATETIME,
  updated_at DATETIME
  PRIMARY KEY (`id`)
);