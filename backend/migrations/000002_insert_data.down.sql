-- Remove all inserted data in reverse order to respect foreign key constraints
DELETE FROM post_privacy;
DELETE FROM notifications;
DELETE FROM group_messages;
DELETE FROM messages;
DELETE FROM event_responses;
DELETE FROM events;
DELETE FROM group_comments;
DELETE FROM group_posts;
DELETE FROM group_invitations;
DELETE FROM group_members;
DELETE FROM groups;
DELETE FROM follow_requests;
DELETE FROM comments;
DELETE FROM posts;
DELETE FROM sessions;
DELETE FROM users;