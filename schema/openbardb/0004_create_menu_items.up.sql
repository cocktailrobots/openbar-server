call dolt_add('.');
call dolt_commit('-m', 'Pre-migration 0004_create_menu_items.up.sql', '--allow-empty');

CREATE TABLE menu_items (
    menu_name_fk varchar(32) NOT NULL COLLATE utf8mb4_0900_ai_ci,
    recipe_id varchar(36) NOT NULL COLLATE utf8mb4_0900_ai_ci,

    FOREIGN KEY (menu_name_fk) REFERENCES menus(name) ON DELETE CASCADE,
    UNIQUE KEY (menu_name_fk, recipe_id)
);

call dolt_add('.');
call dolt_commit('-m', 'Post-migration 0004_create_menu_items.up.sql');
