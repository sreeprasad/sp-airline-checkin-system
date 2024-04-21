#!/bin/bash
set -e

echo "creating schema"

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "mydatabase4" <<-EOSQL
    CREATE TABLE IF NOT EXISTS public.airlines (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL
    );

    CREATE TABLE IF NOT EXISTS public.flights (
        id SERIAL PRIMARY KEY,
        airline_id INTEGER NOT NULL, 
        name VARCHAR(255) NOT NULL,
        CONSTRAINT fk_airline
            FOREIGN KEY (airline_id)
            REFERENCES public.airlines (id)
            ON DELETE CASCADE
    );

    CREATE TABLE IF NOT EXISTS public.trips (
        id SERIAL PRIMARY KEY,
        flight_id INTEGER NOT NULL, 
        flight_time TIMESTAMP NOT NULL,
        CONSTRAINT fk_flight
            FOREIGN KEY (flight_id)
            REFERENCES public.flights (id)
            ON DELETE CASCADE
    );

    CREATE TABLE IF NOT EXISTS public.users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL
    );


    CREATE TABLE public.seats (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        trip_id INTEGER NOT NULL,
        user_id INTEGER, 

        CONSTRAINT fk_trip
            FOREIGN KEY (trip_id)
            REFERENCES public.trips (id)
            ON DELETE CASCADE,
        CONSTRAINT fk_user
            FOREIGN KEY (user_id)
            REFERENCES public.users (id)
            ON DELETE CASCADE
    );

    
EOSQL

