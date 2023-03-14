-- User is a reserved keyword from postgres. 
-- Needs to be inside quotes
CREATE TABLE IF NOT EXISTS users(
    Id SERIAL PRIMARY KEY,
    Username VARCHAR(255),
    Email VARCHAR(255),
    Phone VARCHAR(255),
    PasswordHash VARCHAR(255),
    Roles VARCHAR[],
    CONSTRAINT username_unique UNIQUE (Username)
);