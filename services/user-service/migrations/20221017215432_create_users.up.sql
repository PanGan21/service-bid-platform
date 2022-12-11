-- User is a reserved keyword from postgres. 
-- Needs to be inside quotes
CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    username VARCHAR(255),
    passwordHash VARCHAR(255),
    roles VARCHAR[],
    CONSTRAINT username_unique UNIQUE (username)
);