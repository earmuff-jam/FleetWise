
-- File: 0007_create_storage_locations_table_community.up.sql
-- Description: Create the storage locations table

DROP TABLE IF EXISTS storage_locations CASCADE;
CREATE TABLE IF NOT EXISTS storage_locations
(
    id              UUID PRIMARY KEY                  DEFAULT gen_random_uuid(),
    location        VARCHAR(100)             NOT NULL UNIQUE,
    created_by      UUID REFERENCES profiles (id) ON UPDATE CASCADE ON DELETE CASCADE,
    updated_by      UUID REFERENCES profiles (id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    sharable_groups UUID[]
);

COMMENT ON TABLE storage_locations IS 'location of each storage item belonging to each event';

ALTER TABLE community.storage_locations
    OWNER TO community_admin;

GRANT SELECT, INSERT, UPDATE, DELETE ON community.storage_locations TO community_public;
GRANT SELECT, INSERT, UPDATE, DELETE ON community.storage_locations TO community_test;
GRANT ALL PRIVILEGES ON TABLE community.storage_locations TO community_admin;