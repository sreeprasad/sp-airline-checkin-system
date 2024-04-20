package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bxcodec/faker/v3"
	_ "github.com/lib/pq"
)

type Airline struct {
	ID   uint
	Name string
}

type Flight struct {
	ID        uint
	AirlineID uint
	Name      string
}

type Trip struct {
	ID         uint
	FlightID   uint
	FlightTime time.Time
}

type User struct {
	ID   uint
	Name string
}

type Seat struct {
	ID     uint
	Name   string
	TripID uint
	UserID uint
}

func clearAllContents(db *sql.DB) {
	tables := []string{"seats", "users", "trips", "flights", "airlines"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table))
		if err != nil {
			log.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
	fmt.Println("All tables truncated successfully.")
}

func initializeDB(db *sql.DB) {
	_, err := db.Exec(`INSERT INTO airlines (name) VALUES ('Air India') RETURNING id;`)
	if err != nil {
		log.Fatalf("Failed to insert airline: %v", err)
	}

	var airlineID int
	err = db.QueryRow(`SELECT id FROM airlines WHERE name = 'Air India';`).Scan(&airlineID)
	if err != nil {
		log.Fatalf("Failed to query airline ID: %v", err)
	}

	_, err = db.Exec(`INSERT INTO flights (name, airline_id) VALUES ('AIR_01', $1);`, airlineID)
	if err != nil {
		log.Fatalf("Failed to insert flight: %v", err)
	}

	var flightID int
	err = db.QueryRow(`SELECT id FROM flights WHERE name = 'AIR_01';`).Scan(&flightID)
	if err != nil {
		log.Fatalf("Failed to query flight ID: %v", err)
	}

	specificTime := time.Date(2024, time.April, 19, 21, 0, 0, 0, time.UTC)
	_, err = db.Exec(`INSERT INTO trips (flight_id, flight_time) VALUES ($1, $2);`, flightID, specificTime)
	if err != nil {
		log.Fatalf("Failed to insert trip: %v", err)
	}

	for i := 0; i < 120; i++ {
		userName := faker.Name()
		_, err = db.Exec(`INSERT INTO users (name) VALUES ($1) RETURNING id;`, userName)
		if err != nil {
			log.Fatalf("Failed to insert user: %v", err)
		}

		var userID int
		err = db.QueryRow(`SELECT id FROM users WHERE name = $1;`, userName).Scan(&userID)
		if err != nil {
			log.Fatalf("Failed to query user ID: %v", err)
		}

		seatName := fmt.Sprintf("Seat %d", i+1)
		_, err = db.Exec(`INSERT INTO seats (name, trip_id, user_id) VALUES ($1, $2, $3);`, seatName, flightID, userID)
		if err != nil {
			log.Fatalf("Failed to insert seat: %v", err)
		}
	}

	fmt.Println("Data insertion complete")
}

func ensureAllUsersExist(db *sql.DB) {
	var noNameCount int
	sqlStatement := `SELECT COUNT(*) FROM users WHERE name IS NULL;`
	row := db.QueryRow(sqlStatement)
	err := row.Scan(&noNameCount)
	if err != nil {
		log.Fatalf("Error querying for users without names: %v", err)
	}
	if noNameCount > 0 {
		log.Printf("There are %d users with no name.", noNameCount)
	} else {
		log.Println("All users have names.")
	}
}

func main() {
	connStr := "host=localhost port=5435 user=user4 dbname=mydatabase4 password=password4 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	clearAllContents(db)
	initializeDB(db)
	ensureAllUsersExist(db)
}
