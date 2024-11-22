
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
- https://www.techempower.com/benchmarks

## Data optimizations
- UUID mapping -> int

## JSON optimizations
- Go STD
- third party lib
- save offer as json string in separate varchar column to be able to directly return it in the api response to avoid serialization

## go optimizations:
- GOMAXPROCS env variable
- fasthttp for REST

## sqlite optimizations:
- memory mapping, cache size
- transaktionen!
- batching
- prepared statement
- cached statement

## OS optimizations:
- page size
- file system cache, block size
- tcp: buffer size, timeout, keepalive
- ulimit: open files, max processes
- kernel: max connections, max threads, max memory, max file handles
