-- All passwords are "pass"

INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth, nickname, about_me) VALUES
('amina.ali@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Amina', 'Ali', '1990-05-15', 'Amy', 'Kicks ass');

INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth, nickname, about_me, is_public) VALUES
('yusuf.kamau@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Yusuf', 'Kamau', '1988-11-20', 'Yus', 'Loves hiking and photography.', TRUE);

INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth, nickname, is_public, status) VALUES
('aziz.bek@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Aziz', 'Bek', '1985-07-10', 'AzizB', FALSE, 'enable');

INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth, about_me) VALUES
('zara.khan@example.com', '$2a$10$RhuX9T/SPdJ8vMxGvYcjaufLkseSRvBS5VHwFChHL4W4AwPRlS8bC', 'Zara', 'Khan', '1996-12-25', 'Enjoys reading and learning new languages.');


INSERT INTO posts (user_id, content, privacy_level) VALUES
(1, 'Just enjoyed a beautiful sunset! #nature #peaceful', 'public');

INSERT INTO posts (user_id, content, image_path, privacy_level) VALUES
(2, 'Trying out a new recipe today. Fingers crossed! 🤞', '/images/posts/yusuf_recipe.jpg', 'almost_private');

INSERT INTO posts (user_id, content, privacy_level) VALUES
(3, 'Thinking about life and everything in between.', 'private');

INSERT INTO posts (user_id, content, image_path, privacy_level) VALUES
(4, 'A cozy evening with a good book. 📚', '/images/posts/zara_book.jpg', 'almost_private');

INSERT INTO posts (user_id, content, privacy_level) VALUES
(1, 'Morning coffee and planning the day ahead.', 'public');

INSERT INTO posts (user_id, content, privacy_level) VALUES
(3, 'Some personal reflections I wanted to jot down...', 'private');

INSERT INTO posts (user_id, content, privacy_level) VALUES
(4, 'Found a hidden gem of a cafe! The coffee is amazing. ☕', 'public');