package sp_airlines

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

func ClearAllContents(db *sql.DB) {
	tables := []string{"seats", "users", "trips", "flights", "airlines"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table))
		if err != nil {
			log.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
	fmt.Println("All tables truncated successfully.")
}

func InitializeDB(db *sql.DB) {
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

		seatCode := (i - 1) % 5
		seatRow := seatCode + 1
		seatLetter := rune('A' + (i % 5))
		seatName := fmt.Sprintf("%d-%c", seatRow, seatLetter)

		_, err = db.Exec(`INSERT INTO seats (name, trip_id, user_id) VALUES ($1, $2, $3);`, seatName, flightID, userID)
		if err != nil {
			log.Fatalf("Failed to insert seat: %v", err)
		}
	}

	fmt.Println("Data insertion complete")
}

func ShowAllSeats(db *sql.DB) {
	sqlStatement := `SELECT id, name, trip_id, user_id FROM seats;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	fmt.Println("ID | Name | Trip ID | User ID")
	for rows.Next() {
		var id int
		var name string
		var tripID int
		var userID int
		err = rows.Scan(&id, &name, &tripID, &userID)
		if err != nil {
			log.Fatalf("Failed to read row: %v", err)
		}
		fmt.Printf("%d | %s | %d | %d\n", id, name, tripID, userID)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalf("Error reading rows: %v", err)
	}
}
