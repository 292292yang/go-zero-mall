DROP TABLE IF EXISTS `order`;

CREATE TABLE `order` (
     `id` bigint NOT NULL AUTO_INCREMENT,
     `user_id` bigint NOT NULL,
     `product_id` bigint NOT NULL,
     `quantity` bigint NOT NULL,
     `price` bigint NOT NULL,
     `status` tinyint NOT NULL DEFAULT 1,
     `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
     `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     PRIMARY KEY (`id`),
     KEY `idx_user_id` (`user_id`),
     KEY `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;