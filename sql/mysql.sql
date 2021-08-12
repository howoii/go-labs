CREATE TABLE IF NOT EXISTS account_info
(
    `account_id` BIGINT UNSIGNED NOT NULL ,
    `name` VARCHAR(64) NOT NULL ,
    `create_at` DATETIME DEFAULT NULL,

    PRIMARY KEY(`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

INSERT INTO account_info (account_id, name, create_at) VALUES (10001, 'dijia', CURRENT_TIMESTAMP);