#!/bin/bash
go build -o bookings cmd/web/*.go && ./bookings -dbname=bookings -dbuser=florian.lahitte -dbpass=postgres -dbhost=localhost -dbport=5432 -dbssl=disable