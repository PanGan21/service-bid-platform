CREATE TABLE IF NOT EXISTS requests(
    Id SERIAL PRIMARY KEY,
    Title VARCHAR(255),
    Postcode VARCHAR(255),
    Info VARCHAR(255),
    CreatorId VARCHAR(255),
    Status VARCHAR(255),
    RejectionReason VARCHAR(255)
);