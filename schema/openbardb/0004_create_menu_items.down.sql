call dolt_add('.');
call dolt_commit('-m', 'Pre-migration 0004_create_menu_items.down.sql', '--allow-empty');

DROP TABLE menu_items;

call dolt_add('.');
call dolt_commit('-m', 'Post-migration 0004_create_menu_items.down.sql');
