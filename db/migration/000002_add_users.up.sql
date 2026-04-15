CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL UNIQUE,
  "password_hash" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar NOT NULL UNIQUE,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "password_changed_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "account" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username") DEFERRABLE INITIALLY IMMEDIATE;

-- CREATE INDEX ON "account" ("owner", "currency");
ALTER TABLE "account" ADD CONSTRAINT account_owner_currency_unique UNIQUE ("owner", "currency");