
-- File: 0011_create_category_table.up.sql
-- Description: Create the expenses table

SET search_path TO asset, public;

-- table category is used to stored all categories that an incurred expense can be grouped by --
-- sample category are created by the application. adding new category is allowed --

CREATE TABLE IF NOT EXISTS asset.category
(
    id              UUID PRIMARY KEY             NOT NULL DEFAULT gen_random_uuid(),
    item_name       VARCHAR(100)                 NOT NULL,
    created_by      UUID                         REFERENCES profiles (id) ON UPDATE CASCADE ON DELETE CASCADE,
    updated_by      UUID                         REFERENCES profiles (id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at      TIMESTAMP WITH TIME ZONE     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE     NOT NULL DEFAULT NOW(),
    sharable_groups UUID[]
);

-- prefil category ---
INSERT INTO asset.category (item_name)
VALUES ('Groceries'),
       ('Utilities'),
       ('Rent/Mortgage'),
       ('Entertainment'),
       ('Dining Out'),
       ('Transportation'),
       ('Healthcare'),
       ('Clothing'),
       ('Home Maintenance'),
       ('Education'),
       ('Travel'),
       ('Gifts/Donations'),
       ('Electronics'),
       ('Insurance'),
       ('Personal Care'),
       ('Miscellaneous');

COMMENT ON TABLE category IS 'category is incurred expenses that can be grouped on';

ALTER TABLE asset.category
    OWNER TO asset_admin;

GRANT SELECT, INSERT, UPDATE, DELETE ON asset.category TO asset_public;
GRANT SELECT, INSERT, UPDATE, DELETE ON asset.category TO asset_test;
GRANT ALL PRIVILEGES ON TABLE asset.category TO asset_admin;
