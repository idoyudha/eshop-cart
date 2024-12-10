CREATE TABLE `carts` (
    `id` VARCHAR(36) PRIMARY KEY,
    `user_id` VARCHAR(36) NOT NULL,
    `product_id` VARCHAR(36) NOT NULL,
    `product_name` VARCHAR(255),
    `product_price` FLOAT NOT NULL,
    `product_quantity` INT NOT NULL,
    `note` VARCHAR(255) NOT NULL,
    `updated_at` TIMESTAMP NOT NULL,
    `created_at` TIMESTAMP NOT NULL,
    `deleted_at` TIMESTAMP
);