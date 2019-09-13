CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "public"."whitelists"
(
  "id"          uuid DEFAULT uuid_generate_v4()             PRIMARY KEY NOT NULL,
  "stellar_address"     varchar(200)                        NOT NULL,
  "role"        varchar(100)                                NOT NULL,
  "created_at"  timestamptz                                 NOT NULL,
  "updated_at"  timestamptz                                 NOT NULL,
  "deleted_at"  timestamptz
);

CREATE TABLE "public"."roles"
(
  "id"          serial                                      PRIMARY KEY NOT NULL,
  "name"        varchar(100)                                NOT NULL,
  "role"        varchar(100)                                NOT NULL UNIQUE,
  "created_at"  timestamptz                                 NOT NULL,
  "updated_at"  timestamptz                                 NOT NULL,
  "deleted_at"  timestamptz
);

CREATE UNIQUE INDEX whitelists_partial_unique_code1 ON "public"."whitelists" (address, role, deleted_at)
WHERE deleted_at IS NOT NULL;

CREATE UNIQUE INDEX whitelists_partial_unique_code2 ON "public"."whitelists" (address, role)
WHERE deleted_at IS NULL;

INSERT INTO public.roles
 (name, role, created_at, updated_at) VALUES
('Whitelist user', 'WHITELIST_USER', current_timestamp, current_timestamp),
('Set up token', 'SETUP_TOKEN', current_timestamp, current_timestamp),
('Mint token', 'MINT_TOKEN', current_timestamp, current_timestamp),
('Feed price to Node', 'FEED_PRICE', current_timestamp, current_timestamp);