call dolt_add('.');
call dolt_commit('-m', 'Pre-migration 0001_create_fluids.up.sql', '--allow-empty');

CREATE TABLE fluids (
    idx INT primary key,
    fluid VARCHAR(32)
);

call dolt_add('.');
call dolt_commit('-m', 'Post-migration 0001_create_fluids.up.sql');
