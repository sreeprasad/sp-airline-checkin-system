package main

import (
	airlines "airline-checkin-system/sp_airlines"
	"database/sql"
	"log"
	"sync"

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
	users, err := airlines.GetAllUsers(db)
	if err != nil {
		log.Fatalf("Error in getting all users: %v", err)
		return
	}

	tripID := 1
	var wg sync.WaitGroup
	wg.Add(len(users))
	for idx := range users {
		go func(index int, user *airlines.User) {
			defer wg.Done()

			transaction, err := db.Begin()
			if err != nil {
				log.Printf("Failed to begin transaction: %v", err)
			}

			seat, err := airlines.GetAvailableSeatWithUpdate(transaction, tripID)
			if err != nil {
				log.Fatalf("Invalid input for seat ID: %v", err)
			}

			sqlStatement := `UPDATE seats SET user_id = $1, trip_id = $2 WHERE id = $3;`
			if _, err = transaction.Exec(sqlStatement, user.ID, tripID, seat.ID); err != nil {
				transaction.Rollback()
				log.Printf("Failed to execute update: %v", err)
			}

			if err = transaction.Commit(); err != nil {
				log.Printf("Failed to commit transaction: %v", err)
			}
			//fmt.Printf("%d: User %s was added to seat %s \n", index, user.Name, seat.Name)

		}(idx, &users[idx])
	}
	wg.Wait()

	airlines.PrintAllSeats(db)

}
