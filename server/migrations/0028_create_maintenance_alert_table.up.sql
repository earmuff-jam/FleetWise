
-- File: 0028_create_maintenance_alert_table.up.sql
-- Description: Creates maintenance alert table to alert the users of all the items that are past due in maintenance
-- Note:- updated_by is not forced to use the UUID because the system can make this change automatically --

SET search_path TO community, public;

CREATE TABLE IF NOT EXISTS community.maintenance_alert
(
    id                      UUID                                                                                                         PRIMARY KEY        NOT NULL    DEFAULT gen_random_uuid(),
    maintenance_plan_id     UUID                        UNIQUE REFERENCES maintenance_plan (id) ON UPDATE CASCADE ON DELETE CASCADE,
    name                    VARCHAR(100),
    type                    TEXT                                                                                                                            NOT NULL,
    plan_due                TIMESTAMP WITH TIME ZONE                                                                                                        NOT NULL    DEFAULT NOW(),
    is_read                 BOOLEAN                                                                                                                                     DEFAULT false,
    updated_by              TEXT,
    updated_at              TIMESTAMP WITH TIME ZONE                                                                                                        NOT NULL    DEFAULT NOW(),
    sharable_groups         UUID[]
);

COMMENT ON TABLE maintenance_alert IS 'table to support list of maintenance plans that are due. Due date is derieved from the maintenance plan table.';

ALTER TABLE community.maintenance_alert ADD CONSTRAINT unique_maintenance_plan_id_type UNIQUE (maintenance_plan_id, type);

ALTER TABLE community.maintenance_alert OWNER TO community_admin;

GRANT SELECT, INSERT, UPDATE, DELETE ON community.maintenance_alert TO community_public;
GRANT SELECT, INSERT, UPDATE, DELETE ON community.maintenance_alert TO community_test;
GRANT ALL PRIVILEGES ON TABLE community.maintenance_alert TO community_admin;
