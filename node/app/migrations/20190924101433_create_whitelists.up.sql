CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "public"."whitelists"
(
  "id"          uuid DEFAULT uuid_generate_v4()             PRIMARY KEY NOT NULL,
  "stellar_public_key"     varchar(56)                      NOT NULL,
  "role_code"   varchar(30)                                 NOT NULL,
  "created_at"  timestamptz                                 NOT NULL,
  "updated_at"  timestamptz                                 NOT NULL,
  "deleted_at"  timestamptz
);

CREATE TABLE "public"."roles"
(
  "id"          serial                                      PRIMARY KEY NOT NULL,
  "name"        varchar(100)                                NOT NULL,
  "code"        varchar(30)                                 NOT NULL UNIQUE,
  "created_at"  timestamptz                                 NOT NULL,
  "updated_at"  timestamptz                                 NOT NULL,
  "deleted_at"  timestamptz
);

CREATE UNIQUE INDEX whitelists_partial_unique_code1 ON "public"."whitelists" (stellar_public_key, role_code, deleted_at)
WHERE deleted_at IS NOT NULL;

CREATE UNIQUE INDEX whitelists_partial_unique_code2 ON "public"."whitelists" (stellar_public_key, role_code)
WHERE deleted_at IS NULL;

ALTER TABLE public.whitelists ADD CONSTRAINT whitelists_fk FOREIGN KEY ("role_code") REFERENCES public.roles("code");

INSERT INTO public.roles
 (name, code, created_at, updated_at) VALUES
('Regulator', 'REGULATOR', current_timestamp, current_timestamp),
('Trusted partner', 'TRUSTED_PARTNER', current_timestamp, current_timestamp),
('Price feeder', 'PRICE_FEEDER', current_timestamp, current_timestamp);
