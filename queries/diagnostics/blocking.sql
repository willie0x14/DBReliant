-- Process ID for each PostgreSQL session.
-- application_name identifies the sessions used by the locking labs.
-- state shows whether PostgreSQL considers the session active or idle.
-- wait_event_type and wait_event identify what an active session is waiting on.
-- blocked_by lists the process IDs currently blocking the session.
SELECT
    pid,
    application_name,
    state,
    wait_event_type,
    wait_event,
    pg_blocking_pids(pid) AS blocked_by
FROM pg_stat_activity
WHERE application_name IN (
    'deadlock_a',
    'deadlock_b',
    'migration_writer',
    'migration_writer_1',
    'migration_writer_2',
    'migration_ddl',
    'migration_observer'
);
