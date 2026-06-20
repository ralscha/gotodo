-- +goose Up
CREATE TYPE app_user_authority AS ENUM ('user', 'admin');
CREATE TYPE tokens_scope AS ENUM ('signup', 'password-reset', 'email-change');

CREATE TABLE app_user (
  id            BIGSERIAL PRIMARY KEY,
  email         VARCHAR(255) NOT NULL,
  email_new     VARCHAR(255),
  password_hash VARCHAR(255) NOT NULL,
  authority     app_user_authority NOT NULL,
  activated     BOOLEAN NOT NULL,
  expired       TIMESTAMP,
  last_access   TIMESTAMP,
  UNIQUE(email)
);

CREATE TABLE tokens (
  id          BIGSERIAL PRIMARY KEY,
  hash        BYTEA NOT NULL,
  app_user_id BIGINT NOT NULL REFERENCES app_user(id) ON DELETE CASCADE,
  expiry      TIMESTAMP NOT NULL,
  scope       tokens_scope NOT NULL
);

CREATE TABLE todo (
  id          BIGSERIAL PRIMARY KEY,
  subject     VARCHAR(255) NOT NULL,
  description VARCHAR(255),
  app_user_id BIGINT NOT NULL REFERENCES app_user(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE todo;
DROP TABLE tokens;
DROP TABLE app_user;
DROP TYPE tokens_scope;
DROP TYPE app_user_authority;
