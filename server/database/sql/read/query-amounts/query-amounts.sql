
WITH DataWithFilteredRequiredParams as (
    SELECT * FROM Offers WHERE
        -- regionID_MIN, regionID_MAX
        (mostSpecificRegionID >= ? AND mostSpecificRegionID <= ?) 
        -- timeRangeEnd, timeRangeStart, numberDays
        AND (MAX(endDate, ?) - MIN(startDate, ?) >= ?)
)


-- priceRanges

SELECT 
    FLOOR(price / ?) AS price_range,
    COUNT(*) AS count
FROM DataWithFilteredRequiredParams 
GROUP BY price_range


-- carTypeCounts
-- -- small
-- -- sports
-- -- luxury
-- -- family

SELECT 
    carType,
    COUNT(*) AS count
FROM 
    DataWithFilteredRequiredParams 
GROUP BY 
    carType
ORDER BY 
    carType;


-- seatsCount

SELECT 
    numberSeats,
    COUNT(*) AS count
FROM DataWithFilteredRequiredParams 
GROUP BY numberSeats


-- freeKilometerRange

-- vollkaskoCount
-- -- true
-- -- false