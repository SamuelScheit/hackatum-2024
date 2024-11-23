SELECT json, price
FROM Offers
WHERE
    -- -- -- required -- -- --

    -- regionID_MIN, regionID_MAX
    (mostSpecificRegionID >= ? AND mostSpecificRegionID <= ?) 
    -- timeRangeEnd, timeRangeStart, numberDays
    AND (
        (
            CASE 
                WHEN endDate > ? THEN ?3  
                ELSE endDate
            END
            -
            CASE 
                WHEN startDate < ? THEN ?4 
                ELSE startDate
            END
        ) >= ?
    )
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

LIMIT ? 
OFFSET ?