package main

import (
	"database/sql"
	"log"
	"os"

	"checkmate/types"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var DB_FILE_PATH string = "./db.db"

//go:embed sql/offers-schema.sql
var OFFERS_SCHEMA_SQL string

//go:embed  sql/insert-offer.sql
var INSERT_OFFER_SQL string

func main() {
	initConnection()

	createOffersSchema()
	insertSingleOffer(types.MockOffers[0])

	closeConnection()
}

func initConnection() {
	log.Printf("INIT db connection, using file %s", DB_FILE_PATH)
	os.Remove(DB_FILE_PATH)

	dbconn, err := sql.Open("sqlite3", DB_FILE_PATH)
	if err != nil {
		log.Fatal(err)
	}

	db = dbconn
}

func closeConnection() {
	db.Close()
}

func createOffersSchema() {
	_, err := db.Exec(OFFERS_SCHEMA_SQL)
	if err != nil {
		log.Printf("error creating offers schema :%q\nusing:%s", err, DB_FILE_PATH)
	}
}

func insertSingleOffer(offer types.Offer) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	stmt, err := tx.Prepare(INSERT_OFFER_SQL)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		offer.ID,
		offer.Data,
		offer.MostSpecificRegionID,
		offer.StartDate,
		offer.EndDate,
		offer.NumberSeats,
		offer.Price,
		offer.CarType,
		offer.HasVollkasko,
		offer.FreeKilometers,
	)
	if err != nil {
		_ = tx.Rollback()
		log.Fatalf("Failed to execute statement: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}
}

func retrieveAllOffers() {
	// TODO
}
