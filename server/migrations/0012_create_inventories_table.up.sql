-- File: 0012_create_inventories_table.up.sql
-- Description: create inventory items

SET search_path TO community, public;
CREATE TABLE IF NOT EXISTS community.inventory
(
    id                    UUID PRIMARY KEY         NOT NULL                                         DEFAULT gen_random_uuid(),
    name                  VARCHAR(100)             NOT NULL,
    description           VARCHAR(500),
    price                 DECIMAL(10, 4)                                                            DEFAULT 0.00,
    status                VARCHAR(100)                                                              DEFAULT NULL,
    barcode               VARCHAR(100),
    sku                   VARCHAR(100),
    color                 VARCHAR(10)                                                               DEFAULT NULL,
    quantity              INT                                                                       DEFAULT 1,
    bought_at             VARCHAR(500),
    location              VARCHAR(500),
    storage_location_id                                                                                             UUID REFERENCES storage_locations (id) ON UPDATE CASCADE ON DELETE CASCADE,
    is_returnable         BOOLEAN                  NOT NULL                                         DEFAULT false,
    return_location       VARCHAR(200)             NULL,
    return_datetime       TIMESTAMP WITH TIME ZONE NULL,
    return_notes          VARCHAR(250),
    max_weight            INT                                                                       DEFAULT 0,
    min_weight            INT                                                                       DEFAULT 0,
    max_height            INT                                                                       DEFAULT 0,
    min_height            INT                                                                       DEFAULT 0,
    created_at            TIMESTAMP WITH TIME ZONE NOT NULL                                         DEFAULT NOW(),
    created_by                                                                                                      UUID REFERENCES profiles (id) ON UPDATE CASCADE ON DELETE CASCADE,
    updated_at            TIMESTAMP WITH TIME ZONE NOT NULL                                         DEFAULT NOW(),
    updated_by                                                                                                      UUID REFERENCES profiles (id) ON UPDATE CASCADE ON DELETE CASCADE,
    sharable_groups       UUID[]
);

COMMENT ON TABLE inventory IS 'inventory table is to store all inventory itesm for users';

ALTER TABLE community.inventory
    OWNER TO community_admin;

GRANT SELECT, INSERT, UPDATE, DELETE ON community.inventory TO community_public;
GRANT SELECT, INSERT, UPDATE, DELETE ON community.inventory TO community_test;
GRANT ALL PRIVILEGES ON TABLE community.inventory TO community_admin;

