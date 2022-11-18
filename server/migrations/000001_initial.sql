-- +goose Up
CREATE TABLE app_user (
  id               BIGINT PRIMARY KEY AUTO_INCREMENT,
  email            VARCHAR(255) NOT NULL,
  email_new        VARCHAR(255) NULL,
  password_hash    VARCHAR(255) NOT NULL,
  authority        ENUM('user', 'admin') NOT NULL,
  activated        BOOLEAN NOT NULL,
  expired          TIMESTAMP NULL,
  last_access      TIMESTAMP NULL,
  UNIQUE(email)
);

CREATE TABLE tokens (
  id          BIGINT PRIMARY KEY AUTO_INCREMENT,
  hash        TINYBLOB NOT NULL,
  app_user_id BIGINT NOT NULL,
  expiry      TIMESTAMP NOT NULL,
  scope       ENUM('signup', 'password-reset', 'email-change') NOT NULL,
  FOREIGN KEY (app_user_id) REFERENCES app_user(id) ON DELETE CASCADE
);

CREATE TABLE todo (
  id            BIGINT       PRIMARY KEY AUTO_INCREMENT,
  subject       VARCHAR(255) NOT NULL,
  description   VARCHAR(255) NULL,
  app_user_id   BIGINT       NOT NULL,
  FOREIGN KEY (app_user_id) REFERENCES app_user(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE todo;
DROP TABLE tokens;
DROP TABLE app_user;
