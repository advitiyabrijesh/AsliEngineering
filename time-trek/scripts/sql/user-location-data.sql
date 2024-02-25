create table
    `user_location_data` (
                             `user_location_id` int unsigned not null auto_increment primary key,
                             `user_id` int unsigned not null,
                             `latitude` DECIMAL not null,
                             `longitude` DECIMAL not null,
                             `timestamp` TIMESTAMP not null default CURRENT_TIMESTAMP,
                             FOREIGN KEY (user_id) REFERENCES users(user_id)
);