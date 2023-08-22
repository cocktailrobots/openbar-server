call dolt_add('.');
call dolt_commit('-m', 'Pre-migration 0002_create_pumps.up.sql', '--allow-empty');

CREATE TABLE pumps (
    idx INT primary key,
    ml_per_sec float NOT NULL DEFAULT 0.0
);

call dolt_add('.');
call dolt_commit('-m', 'Post-migration 0002_create_pumps.up.sql');
