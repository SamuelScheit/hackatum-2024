
WITH DataWithFilteredRequiredParams as (
    SELECT * FROM Offers WHERE
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
)


SELECT 
    'price_range' AS GroupingType,
    floor(price / ?) * ?8 AS GroupingValue,
    COUNT(*) AS Count
FROM DataWithFilteredRequiredParams
WHERE ( freeKilometers >= ? )
      and (carType = ?)
      and (numberSeats >= ?)
      and (? IS NULL OR ?12 = false OR hasVollkasko = true)
GROUP BY GroupingValue

UNION ALL

SELECT 
    'carType' AS GroupingType,
    carType AS GroupingValue,
    COUNT(*) AS Count
FROM DataWithFilteredRequiredParams
WHERE
    ( price < ?)
    and ( price >= ? )
    and ( freeKilometers >= ?9 )
    and (numberSeats >= ?11)
GROUP BY carType

UNION ALL

SELECT 
    'numberSeats' AS GroupingType,
    CAST(numberSeats AS VARCHAR) AS GroupingValue,
    COUNT(*) AS Count
FROM DataWithFilteredRequiredParams
WHERE ( price < ?13 )
        and ( price >= ?14 )
        and ( freeKilometers >= ?9 )
        and (carType = ?10)
GROUP BY numberSeats

UNION ALL

SELECT 
    'freeKilometerRange' AS GroupingType,
    floor(freeKilometers / ?) * ?9 AS GroupingValue,
    COUNT(*) AS Count
FROM DataWithFilteredRequiredParams
WHERE 
    ( price < ?13)
    and ( price >= ?14 )
    and (carType = ?10)
    and (numberSeats >= ?11)
GROUP BY GroupingValue

UNION ALL

SELECT 
    'hasVollkasko' AS GroupingType,
    CASE WHEN hasVollkasko THEN 'true' ELSE 'false' END AS GroupingValue,
    COUNT(*) AS Count
FROM DataWithFilteredRequiredParams
WHERE ( price < ?13)
        and ( price >= ?14 )
        and ( freeKilometers >= ?9 )
        and ( carType = ?10 )
        and (numberSeats >= ?11)
GROUP BY hasVollkasko;







-- -- priceRanges

-- SELECT 
--     floor(price / ?) ?1 AS price_range,
--     COUNT(*) AS count
-- FROM DataWithFilteredRequiredParams 
-- GROUP BY price_range


-- -- carTypeCounts
-- -- -- small
-- -- -- sports
-- -- -- luxury
-- -- -- family

-- SELECT 
--     carType,
--     COUNT(*) AS count
-- FROM 
--     DataWithFilteredRequiredParams 
-- GROUP BY 
--     carType
-- ORDER BY 
--     carType;


-- -- seatsCount

-- SELECT 
--     numberSeats,
--     COUNT(*) AS count
-- FROM DataWithFilteredRequiredParams 
-- GROUP BY numberSeats


-- -- freeKilometerRange

-- SELECT 
--     floor(freeKilometers / ?) AS freeKilometerRange,
--     COUNT(*) AS count
-- FROM DataWithFilteredRequiredParams 
-- GROUP BY freeKilometerRange


-- -- vollkaskoCount
-- -- -- true
-- -- -- false

-- SELECT
--     (SELECT 
--     COUNT(*) as count
--     FROM DataWithFilteredRequiredParams
--     WHERE hasVollkasko = true) as trueCount,
--     (SELECT
--     COUNT(*) as count
--     FROM DataWithFilteredRequiredParams
--     WHERE hasVollkasko = false) as falseCount
    