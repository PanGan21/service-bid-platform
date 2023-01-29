CREATE TABLE IF NOT EXISTS requests(
    Id SERIAL PRIMARY KEY,
    Title VARCHAR(255),
    Postcode VARCHAR(255),
    Info VARCHAR(255),
    CreatorId VARCHAR(255),
    Deadline BIGINT,
    Status VARCHAR(255),
    WinningBidId VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS bids(
    Id INTEGER PRIMARY KEY,
    Amount FLOAT,
    CreatorId VARCHAR(255),
    RequestId INTEGER REFERENCES requests (Id)
);