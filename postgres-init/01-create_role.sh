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

   END
    \$\$;
EOSQL


