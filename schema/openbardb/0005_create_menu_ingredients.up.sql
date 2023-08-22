call dolt_add('.');
call dolt_commit('-m', 'Pre-migration 0005_create_menu_ingredients.up.sql', '--allow-empty');

CREATE TABLE menu_ingredients (
    menu_name_fk varchar(32) NOT NULL COLLATE utf8mb4_0900_ai_ci,
    ingredient_name varchar(36) NOT NULL COLLATE utf8mb4_0900_ai_ci,

    FOREIGN KEY (menu_name_fk) REFERENCES menus(name) ON DELETE CASCADE
);

call dolt_add('.');
call dolt_commit('-m', 'Post-migration 0005_create_menu_ingredients.up.sql');
