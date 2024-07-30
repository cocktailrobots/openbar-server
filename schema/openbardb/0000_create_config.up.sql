call dolt_add('.');
call dolt_commit('-m', 'Pre-migration 0000_create_config.up.sql', '--allow-empty');

CREATE TABLE config (
    `key` varchar(64) PRIMARY KEY,
    `value` varchar(255) NOT NULL
);

INSERT INTO config VALUES ('num_pumps', '0'), ('default_volume_ml', '133');

call dolt_add('.');
call dolt_commit('-m', 'Post-migration 0000_create_config.up.sql');
