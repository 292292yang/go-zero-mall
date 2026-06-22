DROP TABLE IF EXISTS `product`;

CREATE TABLE `product` (
   `id` bigint NOT NULL AUTO_INCREMENT,
   `name` varchar(128) NOT NULL DEFAULT '',
   `description` varchar(1024) NOT NULL DEFAULT '',
   `price` bigint NOT NULL DEFAULT 0,
   `stock` bigint NOT NULL DEFAULT 0,
   `status` tinyint NOT NULL DEFAULT 1,
   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
   `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   KEY `idx_status` (`status`),
   KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;