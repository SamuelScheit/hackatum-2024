# Check24 Challenge (Team Checkmate)

[The challenge](https://devpost.com/software/checkmate-fvro9a) was to implement a high-performance backend for a car rental application that enables partners to create offers, customers to search and read offers using REST APIs with filters, optimizing read-write performance.

## How we built it

-   We started with a database and and rest api system
-   we optimized and tested a lot with sql and fasthttp to get every query to at least two digits of milliseconds
-   because of performance issues and the notice on the discord we decided to switch to an in memory search (really late into it at around 23.Nov 24:00)
-   We adopted a highly efficient approach by representing different selection criteria as bit maps, enabling filtering through simple binary AND operations, which significantly enhances speed and concurrency.
-   implemented most of the in memory filters and got it working up to challenge 3 in less than 10 hours

## Challenges we ran into

-   optimizing the sql queries to perform optimally
-   optimizing datastructures

## Accomplishments that we're proud of

-   fully working SQL system
-   really cool (almost fully) working new in memory system with a nice approach to have easy concatenation of queries on the binary level
-   fast and efficient REST API handlers
-   written in go
-   using btree to optimize search

## What we learned

-   a lot of sqlite specialities
-   the go sqlite driver
-   how to code under PREASURE
-   about datastructures and how to use them to optimize performance for dataprocessing and filtering
