-- User is a reserved keyword from postgres. 
-- Needs to be inside quotes
CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY,
    username VARCHAR(255),
    passwordHash VARCHAR(255),
    CONSTRAINT username_unique UNIQUE (username)
);