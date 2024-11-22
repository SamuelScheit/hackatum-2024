DROP TABLE IF EXISTS Offers;
CREATE TABLE Offers (
    ID UUID PRIMARY KEY,
    data TEXT,
    mostSpecificRegionID INT NOT NULL,
    startDate BIGINT NOT NULL,
    endDate BIGINT NOT NULL,  
    numberSeats INT NOT NULL,
    price INT NOT NULL,
    carType VARCHAR(50) NOT NULL,
    hasVollkasko BOOLEAN NOT NULL,
    freeKilometers INT NOT NULL
);
