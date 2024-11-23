package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"checkmate/types"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var DB_FILE_PATH string = "./db.db"

//go:embed sql/schema/offers-schema.sql
var OFFERS_SCHEMA_SQL string

//go:embed sql/create/insert-offer.sql
var INSERT_OFFER_SQL string

//go:embed sql/read/select-offers.sql
var SELECT_ALL_OFFERS_SQL string

//go:embed sql/delete/delete-offers.sql
var DELETE_ALL_OFFERS_SQL string

func test() {
	Init()
	defer CloseConnection()

	InsertSingleOffer(types.MockOffers[0])

	offers, err := RetrieveAllOffers()
	if err != nil {
		fmt.Errorf("Error retrieving offers: %v", err)
	}

	for _, offer := range offers {
		log.Printf("Offer: %+v", offer)
	}
}

func Init() {
	initConnection()
	createOffersSchema()
	initDataQuery()
	initAmountQuery()
	// DeleteAllOffers()
}

func initConnection() {
	log.Printf("INIT db connection, using file %s", DB_FILE_PATH)
	os.Remove(DB_FILE_PATH)

	dbconn, err := sql.Open("sqlite3", DB_FILE_PATH)
	if err != nil {
		panic(err)
	}

	db = dbconn

	db.Exec("")
}

func CloseConnection() {
	db.Close()
}

func createOffersSchema() {
	_, err := db.Exec(OFFERS_SCHEMA_SQL)
	if err != nil {
		log.Printf("error creating offers schema :%q\nusing:%s", err, DB_FILE_PATH)
	}
}

func InsertOffers(offers []types.Offer) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(INSERT_OFFER_SQL)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	searchResultOffer := types.SearchResultOffer{}

	for _, offer := range offers {
		searchResultOffer.Data = offer.Data
		searchResultOffer.ID = offer.ID

		offerJson, err := json.Marshal(searchResultOffer)

		if err != nil {
			fmt.Errorf("failed to marshal offer: %v", err)
			return err
		}

		numberOfDays := offer.EndDate - offer.StartDate

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
			numberOfDays,
			offerJson,
		)

		if err != nil {
			tx.Rollback()
			fmt.Errorf("failed to execute statement: %v", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return err
}

func InsertSingleOffer(offer types.Offer) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	stmt, err := tx.Prepare(INSERT_OFFER_SQL)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	searchResultOffer := types.SearchResultOffer{
		Data: offer.Data,
		ID:   offer.ID,
	}

	offerJson, err := json.Marshal(searchResultOffer)
	if err != nil {
		tx.Rollback()
		return err
	}

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
		offerJson,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute statement: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func RetrieveAllOffers() ([]types.Offer, error) {
	rows, err := db.Query(SELECT_ALL_OFFERS_SQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []types.Offer

	for rows.Next() {
		var offer types.Offer
		err := rows.Scan(
			&offer.ID,
			&offer.Data,
			&offer.MostSpecificRegionID,
			&offer.StartDate,
			&offer.EndDate,
			&offer.NumberSeats,
			&offer.Price,
			&offer.CarType,
			&offer.HasVollkasko,
			&offer.FreeKilometers,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		offers = append(offers, offer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return offers, nil
}

func DeleteAllOffers() error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	_, err = tx.Exec(DELETE_ALL_OFFERS_SQL)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete all offers: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
