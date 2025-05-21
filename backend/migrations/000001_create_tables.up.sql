-- Creating users table to store user information
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    date_of_birth DATE NOT NULL,
    avatar_path TEXT,
    nickname TEXT,
    about_me TEXT,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    status TEXT NOT NULL CHECK (status IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    FOREIGN KEY (updated_by) REFERENCES users (id)
);
-- Creating sessions table for user authentication
CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_token TEXT NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- Creating posts table for user posts
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image_path TEXT,
    privacy_level TEXT NOT NULL CHECK (
        privacy_level IN ('public', 'almost_private', 'private')
    ),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    status TEXT NOT NULL CHECK (status IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- Creating comments table for post comments
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image_path TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    status TEXT NOT NULL CHECK (status IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- Creating follow_requests table for managing follow relationships
CREATE TABLE IF NOT EXISTS follow_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    follower_id INTEGER NOT NULL,
    followed_id INTEGER NOT NULL,
    approval_status TEXT NOT NULL CHECK (
        approval_status IN ('pending', 'accepted', 'declined')
    ),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (followed_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(follower_id, followed_id)
);
-- Creating groups table for group information
CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    creator_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    status TEXT NOT NULL CHECK (status IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
);
-- Creating group_members table for group membership
CREATE TABLE IF NOT EXISTS group_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    approval_status TEXT NOT NULL CHECK (
        approval_status IN ('pending', 'accepted', 'declined')
    ),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    status TEXT NOT NULL CHECK (status IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(group_id, user_id)
);
-- Creating group_invitations table for group membership
CREATE TABLE IF NOT EXISTS group_invitations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    inviter_id INTEGER NOT NULL,
    approval_status TEXT NOT NULL CHECK (
        approval_status IN ('pending', 'accepted', 'declined')
    ),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    status TEXT NOT NULL CHECK (status IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (inviter_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(group_id, user_id)
);
-- Creating group_posts table for group-specific posts
CREATE TABLE IF NOT EXISTS group_posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image_path TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    status TEXT NOT NULL CHECK (status IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- Creating group_comments table for group post comments
CREATE TABLE IF NOT EXISTS group_comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image_path TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    status TEXT NOT NULL CHECK (status IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (group_post_id) REFERENCES group_posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- Creating events table for group events
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    creator_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    event_datetime DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER,
    status TEXT NOT NULL CHECK (status IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
);
-- Creating event_responses table for event attendance
CREATE TABLE IF NOT EXISTS event_responses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    response TEXT NOT NULL CHECK (response IN ('going', 'not_going', 'pending')),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(event_id, user_id)
);
-- Creating messages table for private chats
CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    updated_by INTEGER,
    status TEXT NOT NULL CHECK (status IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
);
-- Creating group_messages table for group chats
CREATE TABLE IF NOT EXISTS group_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    sender_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    updated_by INTEGER,
    status TEXT NOT NULL CHECK (status IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE
);
-- Creating notifications table for user notifications
CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    type TEXT NOT NULL CHECK (
        type IN (
            'follow_request',
            'group_invitation',
            'group_join_request',
            'event_creation'
        )
    ),
    follow_req_id INTEGER,
    group_invite_id INTEGER,
    group_id INTEGER,
    event_id INTEGER,
    content TEXT,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    updated_by INTEGER,
    status TEXT NOT NULL CHECK (
        status IN (
            'enable', 
            'disable', 
            'delete'
        )
    ) DEFAULT 'enable',
    FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (follow_req_id) REFERENCES follow_requests(id) ON DELETE CASCADE,
    FOREIGN KEY (group_invite_id) REFERENCES group_invitations(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);
-- Creating post_privacy table for specific user access to private posts
CREATE TABLE IF NOT EXISTS post_privacy (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    updated_by INTEGER,
    status TEXT NOT NULL CHECK (status IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
    FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(post_id, user_id)
);
-- Creating indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_follow_requests_follower_id ON follow_requests(follower_id);
CREATE INDEX IF NOT EXISTS idx_follow_requests_followed_id ON follow_requests(followed_id);
CREATE INDEX IF NOT EXISTS idx_groups_creator_id ON groups(creator_id);
CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON group_members(group_id);
CREATE INDEX IF NOT EXISTS idx_group_posts_group_id ON group_posts(group_id);
CREATE INDEX IF NOT EXISTS idx_events_group_id ON events(group_id);
CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_receiver_id ON messages(receiver_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);