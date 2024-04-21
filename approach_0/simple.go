package main

import (
	airlines "airline-checkin-system/sp_airlines"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {

	connStr := "host=localhost port=6432 user=user4 dbname=mydatabase4 password=password4 sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	airlines.InitializeDBRecords(db)

	fmt.Printf("enter the user id: ")
	var userID int
	_, err = fmt.Scanln(&userID)
	if err != nil {
		log.Fatalf("Invalid input for user ID: %v", err)
	}

	transaction, err := db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
	}

	user, err := airlines.GetUser(transaction, userID)
	if err != nil {
		log.Fatalf("Invalid input for user ID: %v", err)
	} else {
		fmt.Printf("Welcome %s to SP Airlines\n", user.Name)
	}

	fmt.Printf("enter the seat id: ")
	var seatID int
	_, err = fmt.Scanln(&seatID)
	if err != nil {
		log.Fatalf("Invalid input for seat ID: %v", err)
	}

	seat, err := airlines.GetSeatByID(transaction, seatID)
	if err != nil {
		log.Fatalf("Invalid input for seat ID: %v", err)
	}

	tripID := 1

	sqlStatement := `UPDATE seats SET user_id = $1, trip_id = $2 WHERE id = $3;`
	_, err = transaction.Exec(sqlStatement, seat.ID, tripID, user.ID)

	if err != nil {
		transaction.Rollback()
		fmt.Errorf("execute insert: %w", err)
	}

	if err := transaction.Commit(); err != nil {
		fmt.Errorf("commit transaction: %w", err)
	}

	fmt.Printf("User %s was added to seat %s \n", user.Name, seat.Name)
	airlines.PrettyPrintAllSeats(db)

}
