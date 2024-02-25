create table
    `users` (
                `user_id` int unsigned not null auto_increment primary key,
                `created_at` timestamp not null default CURRENT_TIMESTAMP,
                `email` varchar(255) unique not null,
                `password` varchar(255) null,
                `first_name` varchar(50) null,
                `last_name` varchar(50) null
);