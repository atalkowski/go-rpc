

dft:		# List the targets in this makefile by invoking make
	@echo The default make is to show the targets made by this Makefile and not actually do anything 
	@cat Makefile | grep '^[a-zA-Z0-9_-]*:'

status-postgres:	# Show postgres status but do not change anything.
	postgres.sh

start-postgres:	# Start the postgres DB container.
	postgres.sh start

stop-postgres:	# Stop the postgres DB container.
	postgres.sh stop

db-up:		# Initialise simple_bank db in postgres container 
	docker exec -i postgres createdb --username=root --owner=root simple_bank
	migrate -path  db/migration -database "postgresql://root:mysecret@localhost:5432/simple_bank?sslmode=disable" -verbose up 

db-down:	# Drop the simple_bank db that was set up above.
	docker exec -i postgres dropdb simple_bank

list-schema:	# Display the schema for postgres (db/postgress-accounts.sql) and schema image.
	@cat db/postgres-accounts.sql
	@open images/schema.png

connect:	# Connect to postgres in the docker container.
	@echo Connecting to docker postgres ... use backslash q to quit when done...
	docker exec -it postgres psql

downloads: d1 d2 d3 # The following are the list of installs and downloads.

d1:		# Pull the postgres image from docker.
	./checkit.sh "docker pull postgres:16.2-alpine3.19" "Pull docker image for postgres"

d2:		# Install the psql client (as opposed to TablePlus which costs 90$ or free but very restricted).
	brew install libpg
	@echo "You need to configure your PATH to have the "

d3:		# Install golang migrate to handle database up / down calls - see below.
	brew install golang-migrate

d4:		# Re-initialize the migration scripts used by golang migrate feature in this project.
	rm db/migration/0000*
	migrate create -ext sql -dir db/migration -seq init_schema
	cat db/postgres-accounts.sql > db/migration/000001_init_schema.up.sql
	cat db/postgres-drop.sql     > db/migration/000001_init_schema.down.sql