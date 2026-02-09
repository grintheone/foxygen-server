#!/bin/bash

DB_NAME="foxygendb"
DB_USER="grintheone"

echo "Dropping existing tables..."
psql -U $DB_USER -d $DB_NAME -f ./db/drop_tables.sql

echo "Creating tables..."
psql -U $DB_USER -d $DB_NAME -f ./db/create_tables.sql

echo "Seeding data..."
psql -U $DB_USER -d $DB_NAME -f ./db/seed_data.sql


