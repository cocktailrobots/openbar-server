call dolt_add('.');
call dolt_commit('-m', 'Pre-migration 0003_create_menus.down.sql', '--allow-empty');

DROP TABLE menus;

call dolt_add('.');
call dolt_commit('-m', 'Post-migration 0003_create_menus.down.sql');
