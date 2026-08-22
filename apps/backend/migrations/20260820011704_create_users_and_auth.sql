-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "display_name" text NOT NULL,
  "status" text NOT NULL DEFAULT 'active',
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_users_status" CHECK (status = ANY (ARRAY['active'::text, 'suspended'::text, 'deleted'::text]))
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- Create "auth_identities" table
CREATE TABLE "public"."auth_identities" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "kind" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_auth_identities" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_auth_identity_kind" CHECK (kind = ANY (ARRAY['password'::text, 'google'::text, 'github'::text]))
);
-- Create index "idx_auth_identities_deleted_at" to table: "auth_identities"
CREATE INDEX "idx_auth_identities_deleted_at" ON "public"."auth_identities" ("deleted_at");
-- Create index "uq_auth_identity_user_kind" to table: "auth_identities"
CREATE UNIQUE INDEX "uq_auth_identity_user_kind" ON "public"."auth_identities" ("user_id", "kind");
-- Create "user_emails" table
CREATE TABLE "public"."user_emails" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "email" text NOT NULL,
  "normalized_email" text NOT NULL,
  "is_primary" boolean NOT NULL,
  "verified" boolean NOT NULL DEFAULT false,
  "verified_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_email" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_user_emails_deleted_at" to table: "user_emails"
CREATE INDEX "idx_user_emails_deleted_at" ON "public"."user_emails" ("deleted_at");
-- Create index "idx_user_emails_normalized_email" to table: "user_emails"
CREATE UNIQUE INDEX "idx_user_emails_normalized_email" ON "public"."user_emails" ("normalized_email");
-- Create index "uq_user_emails_primary" to table: "user_emails"
CREATE UNIQUE INDEX "uq_user_emails_primary" ON "public"."user_emails" ("user_id") WHERE (is_primary = true AND deleted_at IS NULL);
