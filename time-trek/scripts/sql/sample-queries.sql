INSERT INTO users (email, password, first_name, last_name)
VALUES ('advitiya@example.com', 'hashed_password', 'Advitiya', 'Brijesh');

INSERT INTO user_location_data (user_id, latitude, longitude)
VALUES (1, 40.7128, -74.0060);

INSERT INTO top_locations (latitude, longitude, type, rating)
VALUES (40.7128, -74.0060, 'Restaurant', 4.5);