CREATE TABLE IF NOT EXISTS requests(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255),
    postcode VARCHAR(255),
    info VARCHAR(255),
    creatorId VARCHAR(255),
    deadline INTEGER
);