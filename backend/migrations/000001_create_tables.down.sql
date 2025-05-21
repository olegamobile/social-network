-- Drop tables in reverse order of creation to avoid foreign key constraints
DROP TABLE IF EXISTS post_privacy;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS group_messages;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS event_responses;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS group_comments;
DROP TABLE IF EXISTS group_posts;
DROP TABLE IF EXISTS group_invitations;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS follow_requests;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
