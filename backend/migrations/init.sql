DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS sessions;

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT,
    email TEXT,
    first_name TEXT,
    last_name TEXT,
    password TEXT,
    birthday TEXT
);

CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    content TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id)
);

INSERT INTO users (username, email, first_name, last_name, password, birthday)
VALUES
('alice', 'alice@example.com', 'Alice', 'Smith', 'pass123', '1990-01-01'),
('bob', 'bob@example.com', 'Bob', 'Jones', 'pass456', '1992-06-15');

INSERT INTO posts (user_id, content)
VALUES
(1, 'Hello, this is Alice''s first post.'),
(1, 'Alice again with another update!'),
(2, 'Bob here. Nice to meet you all.'),
(2, 'Second post from Bob!');
