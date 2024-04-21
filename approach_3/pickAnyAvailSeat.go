package main

import (
	airlines "airline-checkin-system/sp_airlines"
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

func main() {

	connStr := "host=localhost port=5435 user=user4 dbname=mydatabase4 password=password4 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	var tripID = 1

	airlines.InitializeDB(db)

	users, err := airlines.GetAllUsers(db)
	if err != nil {
		log.Fatalf("Error in getting all users: %v", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(users))
	for idx := range users {
		// adding this sleep reduces thread contenion on the seat and all seats
		// will be filled.
		//time.Sleep(10 * time.Millisecond)
		go func(index int, user *airlines.User) {
			defer wg.Done()

			seat, err := airlines.GetAvailableSeat(db, tripID)
			if err != nil {
				log.Fatalf("Invalid input for seat ID: %v", err)
			}
			addTheGodDamUser(db, *user, seat, tripID)
			//fmt.Printf("%d: User %s was added to seat %s \n", index, user.Name, seat.Name)

		}(idx, &users[idx])
	}
	wg.Wait()

	airlines.PrintAllSeats(db)

}

func addTheGodDamUser(db *sql.DB, user airlines.User, seat airlines.Seat, tripID int) error {
	transaction, err := db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return err
	}
	sqlStatement := `UPDATE seats SET user_id = $1, trip_id = $2 WHERE id = $3;`
	if _, err = transaction.Exec(sqlStatement, user.ID, tripID, seat.ID); err != nil {
		transaction.Rollback()
		log.Printf("Failed to execute update: %v", err)
		return err
	}
	if err = transaction.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return err
	}
	return nil
}
