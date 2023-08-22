CREATE DATABASE IF NOT EXISTS cocktails;
USE cocktails;

BEGIN;
CREATE TABLE `cocktails` (
    name VARCHAR(32) primary key NOT NULL COLLATE utf8mb4_0900_ai_ci,
    display_name VARCHAR(64) NOT NULL,
    description TEXT
);

CREATE TRIGGER lowercase_cocktail_name
    BEFORE INSERT ON cocktails
    FOR EACH ROW
        SET NEW.name = LOWER(NEW.name);

CREATE TABLE `ingredients` (
    name VARCHAR(32) primary key NOT NULL COLLATE utf8mb4_0900_ai_ci,
    display_name VARCHAR(64) NOT NULL,
    description TEXT
);

CREATE TRIGGER lowercase_ingredient_name
    BEFORE INSERT ON ingredients
    FOR EACH ROW
        SET NEW.name = LOWER(NEW.name);

CREATE TABLE `recipes` (
    id VARCHAR(36) primary key NOT NULL COLLATE utf8mb4_0900_ai_ci,
    display_name VARCHAR(64) NOT NULL,
    cocktail_fk VARCHAR(32) NOT NULL,
    description TEXT,
    directions TEXT,

    FOREIGN KEY (cocktail_fk) REFERENCES cocktails(name)
);

CREATE TRIGGER lowercase_recipe_id
    BEFORE INSERT ON recipes
    FOR EACH ROW
        SET NEW.id = LOWER(NEW.id);

CREATE TABLE `recipe_ingredients` (
    recipe_id_fk VARCHAR(36) NOT NULL,
    ingredient_fk VARCHAR(32) NOT NULL,
    amount FLOAT NOT NULL,

    PRIMARY KEY (recipe_id_fk, ingredient_fk),
    FOREIGN KEY (recipe_id_fk) REFERENCES recipes(id),
    FOREIGN KEY (ingredient_fk) REFERENCES ingredients(name)
);

COMMIT;
call dolt_commit('-Am', 'Initial schema');