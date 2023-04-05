CREATE TABLE public.todos (
  id VARCHAR(36) PRIMARY KEY,
  content VARCHAR(255) NOT NULL,
  priority SMALLINT CHECK (priority >= 0 AND priority <= 5) NOT NULL DEFAULT 0,
  is_done BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);