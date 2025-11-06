CREATE TABLE IF NOT EXISTS users (
  id            BIGSERIAL PRIMARY KEY,
  first_name    VARCHAR(100) NOT NULL,
  last_name     VARCHAR(100) NOT NULL,
  address       TEXT NOT NULL,
  email         VARCHAR(191) NOT NULL UNIQUE,
  username      VARCHAR(64)  NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  age           INT NOT NULL CHECK (age >= 13),
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT email_format_chk CHECK (
    email ~* '^[A-Z0-9._%+\-]+@[A-Z0-9.\-]+\.[A-Z]{2,}$'
  )
);

CREATE TABLE IF NOT EXISTS categories (
  id   BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS articles (
  id           BIGSERIAL PRIMARY KEY,
  title        VARCHAR(200) NOT NULL,
  content      TEXT NOT NULL,
  author_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  category_id  BIGINT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS likes (
  id          BIGSERIAL PRIMARY KEY,
  user_id     BIGINT NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
  article_id  BIGINT NOT NULL REFERENCES articles(id)  ON DELETE CASCADE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (user_id, article_id)
);

CREATE TABLE IF NOT EXISTS user_activity_logs (
  id          BIGSERIAL PRIMARY KEY,
  user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  action      VARCHAR(64) NOT NULL,        
  description TEXT NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


INSERT INTO categories(name) VALUES
  ('General'), ('Tech'), ('Lifestyle')
ON CONFLICT (name) DO NOTHING;

