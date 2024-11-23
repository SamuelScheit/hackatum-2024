CREATE TABLE IF NOT EXISTS Offers (
    ID varchar(36) PRIMARY KEY,
    data TEXT,
    mostSpecificRegionID INTEGER NOT NULL,
    startDate INTEGER NOT NULL,
    endDate INTEGER NOT NULL,  
    numberSeats INTEGER NOT NULL,
    price INTEGER NOT NULL,
    carType INTEGER NOT NULL,
    hasVollkasko BOOLEAN NOT NULL,
    freeKilometers INTEGER NOT NULL,
    json BLOB NOT NULL
);
CREATE INDEX IF NOT EXISTS mostSpecificRegionID ON Offers (mostSpecificRegionID ASC);
CREATE INDEX IF NOT EXISTS startDate ON Offers (startDate ASC);
CREATE INDEX IF NOT EXISTS endDate ON Offers (endDate ASC);
CREATE INDEX IF NOT EXISTS numberSeats ON Offers (numberSeats);
CREATE INDEX IF NOT EXISTS price ON Offers (price ASC);
CREATE INDEX IF NOT EXISTS price ON Offers (price DESC);
CREATE INDEX IF NOT EXISTS carType ON Offers (carType);
CREATE INDEX IF NOT EXISTS hasVollkasko ON Offers (hasVollkasko);
CREATE INDEX IF NOT EXISTS freeKilometers ON Offers (freeKilometers);
CREATE UNIQUE INDEX IF NOT EXISTS carTypeGroup ON Offers (carType);