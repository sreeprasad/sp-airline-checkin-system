package main

import (
	"fmt"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/jinzhu/gorm"
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

func initializeDB(db *gorm.DB) {
	db.AutoMigrate(&Airline{}, &Flight{}, &Trip{}, &User{}, &Seat{})

	airline := Airline{Name: "Air India"}
	if err := db.Create(&airline).Error; err != nil {
		fmt.Printf("Error creating airline: %v\n", err)
		return
	}

	flight := Flight{Name: "AIR_01", AirlineID: airline.ID}
	if err := db.Create(&flight).Error; err != nil {
		fmt.Printf("Error creating flight: %v\n", err)
		return
	}

	specificTime := time.Date(2024, time.April, 19, 21, 0, 0, 0, time.UTC)
	trip := Trip{FlightID: flight.ID, FlightTime: specificTime}
	if err := db.Create(&trip).Error; err != nil {
		fmt.Printf("Error creating trip: %v\n", err)
		return
	}

	for i := 0; i < 120; i++ {
		user := User{Name: faker.Name()}
		if err := db.Create(&user).Error; err != nil {
			fmt.Printf("Error creating user: %v\n", err)
			return
		}

		seat := Seat{Name: fmt.Sprintf("Seat %d", i+1), TripID: trip.ID, UserID: user.ID}
		if err := db.Create(&seat).Error; err != nil {
			fmt.Printf("Error creating seat: %v\n", err)
			return
		}
	}
	fmt.Println("Data insertion complete")
}

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5435 user=user4 dbname=mydatabase4 password=password4 sslmode=disable")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	clearAllContents(db)
	initializeDB(db)

}

func clearAllContents(db *gorm.DB) {
	tables := []string{"seats", "users", "trips", "flights", "airlines"}

	for _, table := range tables {
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table))
	}

	fmt.Println("All tables truncated successfully.")
}
