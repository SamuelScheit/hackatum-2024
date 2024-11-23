SELECT json, price
FROM Offers
WHERE
    -- -- -- required -- -- --

    -- regionID_MIN, regionID_MAX, 
    -- regionID_MIN2, regionID_MAX2
    (
    (mostSpecificRegionID >= ? AND mostSpecificRegionID <= ?) 
    OR (mostSpecificRegionID >= ? AND mostSpecificRegionID <= ?) 
    )
    -- timeRangeEnd, timeRangeStart, numberDays
    AND ( ? >= endDate )
    AND ( ? <= startDate )
    AND ( numberOfDays = ? )

    -- -- -- optional -- -- -- 
    
    -- minNumberSeats (optional)
    AND (numberSeats >= COALESCE(?, numberSeats))
    -- minPrice (optional)
    AND (price >= COALESCE(?, price))
    -- maxPrice (optional)
    AND (price < COALESCE(?, price))
    -- carType (optional)
    AND (carType = COALESCE(?, carType))
    -- onlyVollkasko (optional)
    AND (? IS NULL OR ?12 = false OR hasVollkasko = true)
    -- minFreeKilometer (optional)
    AND (freeKilometers >= COALESCE(?, freeKilometers))


ORDER BY -- sortorder
    price ASC,
    id ASC

LIMIT ? 
OFFSET ?