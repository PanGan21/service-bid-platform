CREATE TABLE IF NOT EXISTS auctions(
    Id INTEGER PRIMARY KEY,
    Title VARCHAR(255),
    Postcode VARCHAR(255),
    Info VARCHAR(255),
    CreatorId VARCHAR(255),
    Deadline BIGINT,
    Status VARCHAR(255),
    WinningBidId VARCHAR(255),
    RejectionReason VARCHAR(255),
    WinnerId VARCHAR(255),
    WinningAmount FLOAT
);

CREATE TABLE IF NOT EXISTS bids(
    Id SERIAL PRIMARY KEY,
    Amount FLOAT,
    CreatorId VARCHAR(255),
    AuctionId INTEGER REFERENCES auctions (Id)
);