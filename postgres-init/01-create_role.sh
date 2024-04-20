#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" <<-EOSQL
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'user4') THEN
            CREATE ROLE user4 WITH LOGIN PASSWORD 'password4';
        ELSE
            RAISE NOTICE 'Role user4 already exists. Skipping creation.';
        END IF;

        IF NOT EXISTS (
            SELECT 1 
            FROM information_schema.tables 
            WHERE table_schema = 'public' 
            AND table_name = 'posts'
        ) THEN
            CREATE TABLE public.posts (
                id serial PRIMARY KEY,
                title VARCHAR(255),
                author VARCHAR(255)
            );
            RAISE NOTICE 'Table posts created successfully.';
        ELSE
            RAISE NOTICE 'Table posts already exists. Skipping creation.';
        END IF;
    END
    \$\$;
EOSQL


