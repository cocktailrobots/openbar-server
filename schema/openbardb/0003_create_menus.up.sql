call dolt_add('.');
call dolt_commit('-m', 'Pre-migration 0003_create_menus.up.sql', '--allow-empty');

CREATE TABLE menus (
    name varchar(32) PRIMARY KEY NOT NULL COLLATE utf8mb4_0900_ai_ci
);

call dolt_add('.');
call dolt_commit('-m', 'Post-migration 0003_create_menus.up.sql');
