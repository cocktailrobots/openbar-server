call dolt_add('.');
call dolt_commit('-m', 'Pre-migration 0001_create_fluids.down.sql', '--allow-empty');

DROP TABLE fluids;

call dolt_add('.');
call dolt_commit('-m', 'Post-migration 0001_create_fluids.down.sql');
