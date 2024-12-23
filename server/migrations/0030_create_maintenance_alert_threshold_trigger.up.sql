
-- File: 0030_create_maintenance_alert_threshold_trigger.up.sql
-- Description: Creates trigger to populate the maintenance alert if the audience member changes
-- the plan due date to be within the 7 days period.

-- Note:- updated_by is not forced to use the UUID because the system can make this change automatically --


SET search_path TO community, public;

--
-- utility fn created to update maintenance alert table --
-- 
DROP FUNCTION IF EXISTS community.update_maintenance_alert_function_on_plan_due_change_fn() CASCADE;
CREATE FUNCTION community.update_maintenance_alert_function_on_plan_due_change_fn()
    RETURNS trigger AS
$$
BEGIN
INSERT INTO community.maintenance_alert (maintenance_plan_id, name, type, plan_due, is_read, updated_by, updated_at, sharable_groups)
    SELECT 
        mp.id, 
        mp.name, 
        mp.plan_type,
        mp.plan_due, 
        false, 
        NEW.updated_by::TEXT,
        NOW(),
        mp.sharable_groups 
    FROM community.maintenance_plan mp
    WHERE mp.plan_due::DATE BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '7 days'
     ON CONFLICT (maintenance_plan_id) DO UPDATE 
    SET 
        is_read = CASE 
            WHEN community.maintenance_alert.plan_due IS DISTINCT FROM EXCLUDED.plan_due 
            THEN false
            ELSE community.maintenance_alert.is_read
        END,
        plan_due = EXCLUDED.plan_due, 
        updated_at = NOW();

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 
-- trigger used to populate the maintenance alert table --
--
DROP TRIGGER IF EXISTS update_maintenance_alert_function_on_plan_due_change_trigger ON community.maintenance_plan;
CREATE TRIGGER update_maintenance_alert_function_on_plan_due_change_trigger
    AFTER INSERT OR UPDATE
        ON community.maintenance_plan
    FOR EACH ROW
EXECUTE FUNCTION community.update_maintenance_alert_function_on_plan_due_change_fn();
