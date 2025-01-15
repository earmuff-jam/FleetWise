-- File: 0031_update_user_verification_table.up.sql
-- Description: Update oauth table to include user verification check

SET search_path TO auth, public;

ALTER TABLE auth.users
ADD COLUMN is_verified BOOLEAN DEFAULT false;