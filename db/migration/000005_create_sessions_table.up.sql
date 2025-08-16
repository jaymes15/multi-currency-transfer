-- Create sessions table
CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- Add foreign key constraint
ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

-- Create index on username for better performance
CREATE INDEX "idx_sessions_username" ON "sessions" ("username");

-- Create index on refresh_token for token lookups
CREATE INDEX "idx_sessions_refresh_token" ON "sessions" ("refresh_token");

-- Create index on expires_at for cleanup operations
CREATE INDEX "idx_sessions_expires_at" ON "sessions" ("expires_at"); 