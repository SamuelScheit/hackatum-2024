SELECT json, price
FROM Offers
WHERE
    -- -- -- required -- -- --

    -- regionID_MIN, regionID_MAX
    (mostSpecificRegionID >= ? AND mostSpecificRegionID <= ?) 
    -- timeRangeEnd, timeRangeStart, numberDays
    AND (MAX(endDate, ?) - MIN(startDate, ?) >= ?)
    -- pagination: previousPrice
    AND (price > ?)

    -- -- -- optional -- -- -- 
    
    -- minNumberSeats (optional)
    AND (numberSeats >= COALESCE(?, numberSeats))
    -- minPrice (optional)
    AND (price >= COALESCE(?, price))
    -- maxPrice (optional)
    AND (price <= COALESCE(?, price))
    -- carType (optional)
    AND (carType = COALESCE(?, carType))
    -- onlyVollkasko (optional)
    AND (hasVollkasko = COALESCE(?, hasVollkasko))
    -- minFreeKilometer (optional)
    AND (freeKilometers >= COALESCE(?, freeKilometers))


ORDER BY -- sortorder
    price ASC

LIMIT ? -- pagesize