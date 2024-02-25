create table
    `top_locations` (
                        `location_id` int unsigned not null auto_increment primary key,
                        `created_at` timestamp not null default CURRENT_TIMESTAMP,
                        `latitude` DECIMAL(10, 6) not null,
                        `longitude` DECIMAL(10, 6) not null,
                        `type` varchar(50) null,
                        `rating` DECIMAL(3, 2) null default 0
);