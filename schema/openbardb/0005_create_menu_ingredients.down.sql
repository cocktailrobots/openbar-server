call dolt_add('.');
call dolt_commit('-m', 'Pre-migration 0005_create_menu_ingredients.down.sql', '--allow-empty');

DROP TABLE menu_ingredients;

call dolt_add('.');
call dolt_commit('-m', 'Post-migration 0005_create_menu_ingredients.down.sql');
