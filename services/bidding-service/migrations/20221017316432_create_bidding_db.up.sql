CREATE TABLE IF NOT EXISTS requests(
    Id INTEGER PRIMARY KEY,
    Title VARCHAR(255),
    Postcode VARCHAR(255),
    Info VARCHAR(255),
    CreatorId VARCHAR(255),
    Deadline INTEGER,
    Status VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS bids(
    Id SERIAL PRIMARY KEY,
    Amount FLOAT,
    CreatorId VARCHAR(255),
    RequestId INTEGER REFERENCES requests (Id)
);