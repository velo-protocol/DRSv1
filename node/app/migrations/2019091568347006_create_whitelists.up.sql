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

CREATE UNIQUE INDEX whitelists_partial_unique_code1 ON "public"."whitelists" (stellar_address, role, deleted_at)
WHERE deleted_at IS NOT NULL;

CREATE UNIQUE INDEX whitelists_partial_unique_code2 ON "public"."whitelists" (stellar_address, role)
WHERE deleted_at IS NULL;

INSERT INTO public.roles
 (name, role, created_at, updated_at) VALUES
('KYC checker', 'REGULATOR', current_timestamp, current_timestamp),
('Trusted partner', 'TRUSTED_PARTNER', current_timestamp, current_timestamp),
('Price feeder', 'PRICE_FEEDER', current_timestamp, current_timestamp),