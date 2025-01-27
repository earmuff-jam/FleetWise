
-- File: 0029_create_maintenance_alert_threshold_fn.up
-- Description: Function used to populate the maintenance alert. If the plan_due date is within the time
-- frame of 7 days, then we want to populate the maintenance alert table.

-- Note:- updated_by is not forced to use the UUID because the system can make this change automatically --


SET search_path TO community, public;

--
-- Utility function that is used to populate maintenance alerts if they pass the threshold of 7 days. --
-- All maintenance plans that cross 7 days mark are displayed in the table --
--

CREATE OR REPLACE FUNCTION community.populate_maintenance_alerts()
RETURNS void AS
$$
BEGIN
    INSERT INTO community.maintenance_alert (maintenance_plan_id, name, type, plan_due, is_read, updated_by, updated_at, sharable_groups)
    SELECT 
        mp.id, 
        mp.name, 
        mp.plan_type, 
        mp.plan_due, 
        false, 
        'system',
        NOW(),
        mp.sharable_groups
    FROM community.maintenance_plan mp
    WHERE mp.plan_due::DATE BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '7 days' 
     ON CONFLICT (maintenance_plan_id) DO UPDATE 
    SET 
        plan_due = EXCLUDED.plan_due, 
        updated_at = NOW();
END;
$$ LANGUAGE plpgsql;


-- Cron Scheduler to run the function populate_maintenance_alerts every day --
-- SELECT cron.schedule('0 0 * * *', $$CALL community.populate_maintenance_alerts();$$);

-- Ensure pg_cron is only used in the correct database
DO $$ 
DECLARE db_name TEXT;
BEGIN
    SELECT current_database() INTO db_name;
    IF db_name = (SELECT setting FROM pg_settings WHERE name = 'cron.database_name') THEN
        -- Schedule the job only if we are in the correct pg_cron database
        PERFORM cron.schedule('daily_maintenance_alerts', '0 0 * * *', 'SELECT community.populate_maintenance_alerts();');
    ELSE
        RAISE NOTICE 'Skipping cron job creation: not in pg_cron database';
    END IF;
END $$;
