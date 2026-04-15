ALTER TABLE "account" DROP CONSTRAINT IF EXISTS account_owner_currency_unique;
ALTER TABLE "account" DROP CONSTRAINT IF EXISTS account_owner_fkey;

drop table if exists "users" cascade;
