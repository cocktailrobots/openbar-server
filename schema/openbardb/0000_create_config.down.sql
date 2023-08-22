call dolt_add('.');
call dolt_commit('-m', 'Pre-migration 0000_create_config.down.sql', '--allow-empty');

DROP TABLE config;

call dolt_add('.');
call dolt_commit('-m', 'Post-migration 0000_create_config.down.sql');

