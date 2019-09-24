CREATE TABLE "public"."price_entries"
(
  "id"                                  SERIAL            PRIMARY KEY NOT NULL,
  "stellar_public_key"                  varchar(56)       NOT NULL,
  "asset"                               varchar(12)       NOT NULL,
  "currency"                            varchar(12)       NOT NULL,
  "price_in_currency_per_asset_unit"    numeric(29,7)     NOT NULL,
  "created_at"                          timestamptz
);
