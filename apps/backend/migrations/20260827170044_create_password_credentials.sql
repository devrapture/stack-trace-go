-- Create "password_credentials" table
CREATE TABLE "public"."password_credentials" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "auth_identity_id" uuid NOT NULL,
  "password_hash" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_password_credentials_auth_identity" FOREIGN KEY ("auth_identity_id") REFERENCES "public"."auth_identities" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_password_credentials_auth_identity_id" to table: "password_credentials"
CREATE UNIQUE INDEX "idx_password_credentials_auth_identity_id" ON "public"."password_credentials" ("auth_identity_id");
-- Create index "idx_password_credentials_deleted_at" to table: "password_credentials"
CREATE INDEX "idx_password_credentials_deleted_at" ON "public"."password_credentials" ("deleted_at");
