DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `username` varchar(64) NOT NULL DEFAULT '',
    `password` varchar(128) NOT NULL DEFAULT '',
    `nickname` varchar(64) NOT NULL DEFAULT '',
    `avatar` varchar(255) NOT NULL DEFAULT '',
    `status` tinyint NOT NULL DEFAULT 1,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;