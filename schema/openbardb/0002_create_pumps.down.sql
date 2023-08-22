call dolt_add('.');
call dolt_commit('-m', 'Pre-migration 0002_create_pumps.down.sql', '--allow-empty');

DROP TABLE pumps;

call dolt_add('.');
call dolt_commit('-m', 'Post-migration 0002_create_pumps.down.sql');
