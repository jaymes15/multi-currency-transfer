-- Create users table
CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "is_email_verified" boolean NOT NULL DEFAULT false,
  "password_changed_at" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),  
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Add foreign key constraint from accounts.owner to users.username
ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

-- Add unique constraint on accounts (owner, currency)
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");

-- Create index on users table for performance
CREATE INDEX ON "users" ("email");
CREATE INDEX ON "users" ("username");

-- Add comments for documentation
COMMENT ON TABLE "users" IS 'User authentication and profile information';
COMMENT ON COLUMN "users"."username" IS 'Unique username for login and account ownership';
COMMENT ON COLUMN "users"."hashed_password" IS 'BCrypt hashed password for security';
COMMENT ON COLUMN "users"."full_name" IS 'User full name for display';
COMMENT ON COLUMN "users"."email" IS 'Unique email address for notifications';
COMMENT ON COLUMN "users"."is_email_verified" IS 'Email verification status for security';
COMMENT ON COLUMN "users"."password_changed_at" IS 'Timestamp of last password change for security';
COMMENT ON COLUMN "users"."created_at" IS 'User account creation timestamp';
COMMENT ON COLUMN "accounts"."owner" IS 'Username reference to users table'; 