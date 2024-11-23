
## Notes from intro meeting
- ubuntu server with root
- interface must be http
- 4 GiB not enough to go for pure in-memory implementation
- https://hackatum.check24.de/docs/main-challenge/introduction
- sidechallenge nice frontend application (has to be webapplication)
- disable default body limit


# Ablauf

incoming POST -> in queue
queue -> batchweise in SQL
 - dabei search trees aufbauen
GET Request -> return results as soon as possible

# Optimizations

Overhead reduzieren:
- keinen reverse proxy
- kein Docker

Ergebnisse schnell:
- Suchbaum aufbauen

- Deserialisierung der eingehenden json daten in go structs


## Search data structure
- AVL
- Red/Black tree
- R*-tree (https://www.sqlite.org/rtree.html): Benchmark builtin sqlite rtree vs custom go implementation (concurrency)

## HTTP server:
- https://github.com/valyala/fasthttp
    - don't forget profiling
	- GOMAXPROCS env variable
	- https://github.com/valyala/fasthttp#performance-optimization-tips-for-multi-core-systems
- https://www.techempower.com/benchmarks

## Data optimizations
- UUID mapping -> int

## JSON optimizations
- Go STD
- third party lib
- save offer as json string in separate varchar column to be able to directly return it in the api response to avoid serialization

## go optimizations:

## sqlite optimizations:
- memory mapping, cache size
- transaktionen!
- batching
- prepared statement
- cached statement
- json1 statements

## OS optimizations:
- page size
- file system cache, block size
- tcp: buffer size, timeout, keepalive
- ulimit: open files, max processes
- kernel: max connections, max threads, max memory, max file handles


# Suchparameteroptimierung

## Regions

- O(1) Dictionarylookup
- mappen eine region auf alle leaf regions die drin sind
- + eventuell min max der leafs

## Timerangestart, Timerangeend, numberOfDays

- Idee: Schnittmenge von zwei tages abschnitten
- R*Star
- z.start_date <= p.query_end AND z.end_date >= p.query_start

## sortOrder

- price index, sql orderby

## page, pageSize

- https://stackoverflow.com/questions/109232/what-is-the-best-way-to-paginate-results-in-sql-server
- seek method / keyset pagination

## praceRangeWidth

```
SELECT 
    FLOOR(price / 500) * 500 AS price_range_start,
    FLOOR(price / 500) * 500 + 499 AS price_range_end,
    COUNT(*) AS num_items
FROM 
    products
GROUP BY 
    FLOOR(price / 500)
ORDER BY 
    price_range_start;
```

## minFreeKilometerWidth

- same as pricerangewidth

## minNumberSeats, minPrice, maxPrice, minFreeKilometer

- easy '>' comparison

## carType

- use integer for cartype and convert to