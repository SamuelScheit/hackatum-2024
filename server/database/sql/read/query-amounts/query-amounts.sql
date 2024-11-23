
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
WHERE ( freeKilometers >= COALESCE(?, freeKilometers  ))
      and (carType = COALESCE(?, carType))
      and (numberSeats >= COALESCE(?, numberSeats))
      and (? IS NULL OR ?12 = false OR hasVollkasko = true)
GROUP BY GroupingValue

UNION ALL

SELECT 
    'carType' AS GroupingType,
    carType AS GroupingValue,
    COUNT(*) AS Count
FROM DataWithFilteredRequiredParams
WHERE
    CASE WHEN ? IS NULL THEN true ELSE ( price < ?13) END 
    and (price >= COALESCE(?, price))
    and (freeKilometers >= COALESCE(?9, freeKilometers))
    and (numberSeats >= COALESCE(?11, numberSeats))
GROUP BY carType

UNION ALL

SELECT 
    'numberSeats' AS GroupingType,
    CAST(numberSeats AS VARCHAR) AS GroupingValue,
    COUNT(*) AS Count
FROM DataWithFilteredRequiredParams
WHERE CASE WHEN ?13 IS NULL THEN true ELSE ( price < ?13) END 
        and (price >= COALESCE(?14, price)) 
        and (freeKilometers >= COALESCE(?9, freeKilometers))
        and (carType = COALESCE(?10, carType))
GROUP BY numberSeats

UNION ALL

SELECT 
    'freeKilometerRange' AS GroupingType,
    floor(freeKilometers / ?) * ?9 AS GroupingValue,
    COUNT(*) AS Count
FROM DataWithFilteredRequiredParams
WHERE 
    CASE WHEN ?13 IS NULL THEN true ELSE ( price < ?13) END 
    and (price >= COALESCE(?14, price)) 
    and (carType = COALESCE(?10, carType))
    and (numberSeats >= COALESCE(?11, numberSeats))
GROUP BY GroupingValue

UNION ALL

SELECT 
    'hasVollkasko' AS GroupingType,
    CASE WHEN hasVollkasko THEN 'true' ELSE 'false' END AS GroupingValue,
    COUNT(*) AS Count
FROM DataWithFilteredRequiredParams
WHERE 
    CASE WHEN ?13 IS NULL THEN true ELSE ( price < ?13) END 
    and (price >= COALESCE(?14, price)) 
    and (freeKilometers >= COALESCE(?9, freeKilometers))
    and (carType = COALESCE(?10, carType))
    and (numberSeats >= COALESCE(?11, numberSeats))
GROUP BY hasVollkasko;

